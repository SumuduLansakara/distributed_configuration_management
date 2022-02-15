#pragma once

#include "Component.h"
#include "RemoteComponent.h"
#include "RemoteComponentStub.h"
#include "Utils.h"
#include <shared_mutex>

namespace cfgStore {
    class LocalComponent : public Component {
    public:
        LocalComponent(std::string const& kind, std::string const& name)
            : Component{kind, name, utils::generateUuidV4()} {
        }

        ~LocalComponent() override;

        void setParam(std::string const& key, std::string const& val) override;
        std::string getParam(std::string const& key) override;
        void removeParam(std::string const& key) override;

        void reloadParam(std::string const& key);
        void reloadParams();

        void initConfigs(std::map<std::string, std::string> configs);

        void connect();
        void disconnect();

        static std::map<std::string, std::vector<std::string>> listComponents(std::string const& kind = "");

        void watchComponents(std::string const& kind);
        void watchParameters(RemoteComponentStub stub);

        virtual void onComponentConnected(RemoteComponentStub stub) = 0;
        virtual void onComponentDisconnected(RemoteComponentStub stub) = 0;

        virtual void onParameterSet(RemoteComponentStub stub, std::string const& key) = 0;
        virtual void onParameterRemoved(RemoteComponentStub stub, std::string const& key) = 0;

    private:
        void writeConfigs();
        void handleComponentEvent(etcd::Response const&);
        void handleParameterEvent(etcd::Response const&);

        bool _connected = false;
        std::map<std::string, std::string> _configs;
        std::vector<std::shared_ptr<etcd::Watcher>> _componentWatchers;
        std::map<std::string, std::pair<RemoteComponentStub, std::shared_ptr<etcd::Watcher>>> _parameterWatchers;

        std::shared_mutex _connectionStateLock;
        std::mutex _paramWatcherLock;
    };
}// namespace cfgStore
