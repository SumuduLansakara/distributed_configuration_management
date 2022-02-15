#include "LocalComponent.h"
#include "Communicator.h"
#include "RemoteComponent.h"
#include <plog/Log.h>

namespace cfgStore {
    LocalComponent::~LocalComponent() {
        PLOGD << "Destructing localComponent: " << prettyName();
    }

    void LocalComponent::setParam(std::string const& key, std::string const& val) {
        Communicator::set(paramPath(_instanceId, key), val);
        _configs[key] = val;
    }

    std::string LocalComponent::getParam(std::string const& key) {
        if (auto const ite = _configs.find(key); ite != _configs.end()) {
            return ite->second;
        }
        throw std::runtime_error(fmt::format("Parameter does not exist {}", key));
    }

    void LocalComponent::removeParam(std::string const& key) {
        Communicator::rm(paramPath(_instanceId, key));
        _configs.erase(key);
    }

    void LocalComponent::reloadParam(const std::string& key) {
        _configs[key] = Communicator::get(paramPath(_instanceId, key));
    }

    void LocalComponent::reloadParams() {
        auto const prefix = paramPath(_instanceId) + '/';
        auto const res = Communicator::ls(prefix);
        for (int i = 0; i < res.keys().size(); ++i) {
            auto const key = res.key(i).substr(prefix.size());
            auto const val = res.value(i).as_string();
            _configs.at(key) = val;
        }
    }

    void LocalComponent::initConfigs(std::map<std::string, std::string> configs) {
        _configs = std::move(configs);
    }

    void LocalComponent::connect() {
        std::unique_lock<std::shared_mutex> g{_connectionStateLock};
        PLOGD << "connecting " << prettyName();
        // parameters
        writeConfigs();

        auto parallelCtx = Communicator::getParallelContext();
        // metadata
        parallelCtx.set(metadataPath(_instanceId, "name"), _name);
        parallelCtx.set(metadataPath(_instanceId, "kind"), _kind);
        // component
        parallelCtx.set(componentPath(_kind, _instanceId), "");
        _connected = true;
    }

    void LocalComponent::disconnect() {
        PLOGD << "disconnecting " << prettyName();
        std::unique_lock<std::shared_mutex> g{_connectionStateLock};
        _connected = true;
        for (auto& watch : _componentWatchers) watch->Cancel();

        {
            std::lock_guard<std::mutex> paramWatchGuard{_paramWatcherLock};
            for (auto& [k, v] : _parameterWatchers) v.second->Cancel();
        }

        auto parallelCtx = Communicator::getParallelContext();
        parallelCtx.rmdir(componentPath(_kind, _instanceId), "");
        parallelCtx.rmdir(metadataPath(_instanceId), _name);
        parallelCtx.rmdir(paramPath(_instanceId), _name);
    }

    std::map<std::string, std::vector<std::string>> LocalComponent::listComponents(std::string const& kind) {
        etcd::Response res;
        auto const basePrefix = componentPath() + '/';
        if (kind.empty()) {
            Communicator::ls(basePrefix);
        } else {
            Communicator::ls(componentPath(kind) + '/');
        }
        std::map<std::string, std::vector<std::string>> ret;
        for (int i = 0; i < res.keys().size(); ++i) {
            auto const key = res.key(i).substr(basePrefix.size());
            auto const tokens = utils::split(key, '/');
            auto const& compKind = tokens.at(0);
            auto const& compInstanceId = tokens.at(1);
            if (auto ite = ret.lower_bound(compKind); ite != ret.end() && ite->first == compKind) {
                ite->second.push_back(compInstanceId);
            } else {
                ret.insert(ite, {compKind, std::vector{compInstanceId}});
            }
        }
        return ret;
    }

    void LocalComponent::watchComponents(std::string const& kind) {
        std::shared_lock<std::shared_mutex> g{_connectionStateLock};
        if (!_connected) {
            PLOGW << "Can't add component watch while disconnected";
            return;
        }
        auto const prefix = componentPath(kind);
        PLOGD << "Watching components: " << prefix;
        auto w = Communicator::createWatcher(
                prefix,
                [this](const etcd::Response& r) { handleComponentEvent(r); },
                true);
        _componentWatchers.emplace_back(std::move(w));
    }

    void LocalComponent::watchParameters(RemoteComponentStub stub) {
        std::shared_lock<std::shared_mutex> g1{_connectionStateLock};
        if (!_connected) {
            PLOGW << "Can't add parameter watch while disconnected";
            return;
        }
        auto const prefix = paramPath(stub.getInstanceId());
        PLOGD << "Watching parameters: " << prefix;
        auto w = Communicator::createWatcher(
                prefix, [this](const etcd::Response& r) { handleParameterEvent(r); },
                true);
        {
            std::lock_guard<std::mutex> g2{_paramWatcherLock};
            _parameterWatchers.emplace(stub.getInstanceId(), std::pair{stub, std::move(w)});
        }
    }

    void LocalComponent::writeConfigs() {
        auto parallelCtx = Communicator::getParallelContext(_configs.size());
        for (auto& [key, val] : _configs) {
            parallelCtx.set(paramPath(_instanceId, key), val);
        }
    }

    void LocalComponent::handleComponentEvent(const etcd::Response& response) {
        auto const prefix = componentPath() + '/';
        for (auto const& event : response.events()) {
            switch (event.event_type()) {
                case etcd::Event::EventType::PUT: {
                    auto const tokens = utils::split(event.kv().key().substr(prefix.size()), '/');
                    PLOGD << "remote component connection detected: " << tokens[0] << " " << tokens[1];
                    onComponentConnected(RemoteComponentStub(tokens[0], tokens[1]));
                    break;
                }
                case etcd::Event::EventType::DELETE_: {
                    auto const tokens = utils::split(event.kv().key().substr(prefix.size()), '/');
                    PLOGD << "remote component disconnection detected: " << tokens[0] << " " << tokens[1];
                    {
                        std::lock_guard<std::mutex> g{_paramWatcherLock};
                        if (auto ite = _parameterWatchers.find(tokens[1]); ite != _parameterWatchers.end()) {
                            PLOGD << "removing parameter watchers of disconnected component: "
                                  << " " << tokens[1];
                            ite->second.second->Cancel();
                        }
                    }
                    onComponentDisconnected(RemoteComponentStub(tokens[0], tokens[1]));
                    break;
                }
                default:
                    PLOGW << "Unsupported component event detected: " << static_cast<int>(event.event_type());
                    break;
            }
        }
    }

    void LocalComponent::handleParameterEvent(const etcd::Response& response) {
        auto const prefix = paramPath() + '/';
        for (auto const& event : response.events()) {
            switch (event.event_type()) {
                case etcd::Event::EventType::PUT: {
                    auto const tokens = utils::split(event.kv().key().substr(prefix.size()), '/');
                    PLOGD << "parameter set detected: " << tokens[0] << " " << tokens[1];
                    {
                        std::lock_guard<std::mutex> g{_paramWatcherLock};
                        if (auto ite = _parameterWatchers.find(tokens[0]); ite != _parameterWatchers.end()) {
                            onParameterSet(ite->second.first, tokens[1]);
                        } else {
                            PLOGD << "ignoring parameter set event of removed watcher: " << tokens[0];
                        }
                    }
                    break;
                }
                case etcd::Event::EventType::DELETE_: {
                    auto const tokens = utils::split(event.kv().key().substr(prefix.size()), '/');
                    PLOGD << "parameter removal detected: " << tokens[0] << " " << tokens[1];
                    // Do not hold `_paramWatcherLock` here, it can cause a deadlock
                    if (auto ite = _parameterWatchers.find(tokens[0]); ite != _parameterWatchers.end()) {
                        onParameterRemoved(ite->second.first, tokens[1]);
                    } else {
                        PLOGD << "ignoring parameter remove event of removed watcher: " << tokens[0];
                    }
                    break;
                }
                default:
                    PLOGW << "Unsupported parameter event detected: " << static_cast<int>(event.event_type());
                    break;
            }
        }
    }
}// namespace cfgStore
