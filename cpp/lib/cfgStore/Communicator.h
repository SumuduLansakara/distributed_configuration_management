#pragma once

#include <etcd/Watcher.hpp>
#include <sstream>
#include <utility>

namespace cfgStore {

    namespace internal {
        template<typename T>
        void concat(std::stringstream& ss, T&& t) {
            ss << t;
        }

        template<typename T, typename... Args>
        void concat(std::stringstream& ss, T&& f, Args... args) {
            ss << f << '/';
            concat(ss, args...);
        }
    }// namespace internal

    template<typename... Args>
    std::string path(Args... args) {
        std::stringstream ss;
        internal::concat(ss, args...);
        return ss.str();
    }

    template<typename... Args>
    std::string paramPath(Args... args) {
        return path(PREFIX_PARAMETERS, args...);
    }

    template<typename... Args>
    std::string metadataPath(Args... args) {
        return path(PREFIX_METADATA, args...);
    }

    template<typename... Args>
    std::string componentPath(Args... args) {
        return path(PREFIX_COMPONENTS, args...);
    }

    class Communicator {
        static inline Communicator* _instance = nullptr;

        class ParallelContext {
        public:
            ParallelContext(etcd::Client& client, size_t taskCount)
                : _client{client}
                , tasks{} {
                tasks.reserve(taskCount);
            }

            ~ParallelContext() {
                waitForCompletion();
            }

            void set(std::string const& key, std::string const& val) {
                tasks.emplace_back(_client.set(key, val));
            }

            void rmdir(std::string const& key, std::string const& val) {
                tasks.emplace_back(_client.rmdir(key));
            }

            void waitForCompletion() {
                for (auto& task : tasks) task.get();
                tasks.clear();
            }

        private:
            std::vector<pplx::task<etcd::Response>> tasks;
            etcd::Client& _client;
        };

    public:
        static Communicator* init(std::string const& etcdEndpoint) {
            if (!_instance) _instance = new Communicator{etcdEndpoint};
            return _instance;
        }

        static void deInit() {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            delete _instance;
            _instance = nullptr;
        }

        static std::shared_ptr<etcd::Watcher> createWatcher(std::string const& key, std::function<void(etcd::Response const&)> callback, bool recursive) {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            return std::make_shared<etcd::Watcher>(_instance->_etcdEndpoint, key, std::move(callback), recursive);
        }

        static void set(std::string const& key, std::string const& val) {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            _instance->_client.set(key, val).get();
        }

        static std::string get(std::string const& key) {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            auto const res = _instance->_client.get(key).get();
            if (!res.is_ok())
                throw std::runtime_error(res.error_message());
            return res.value().as_string();
        }

        static void rm(std::string const& key) {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            _instance->_client.rm(key).get();
        }

        static void rmdir(std::string const& key) {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            _instance->_client.rmdir(key).get();
        }

        static etcd::Response ls(std::string const& prefix) {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            return _instance->_client.ls(prefix).get();
        }

        static ParallelContext getParallelContext(size_t taskCount = 0) {
            if (!_instance) throw std::logic_error("communicator is not initialized");
            return ParallelContext{_instance->_client, taskCount};
        }

    private:
        explicit Communicator(std::string etcdEndpoint)
            : _etcdEndpoint{std::move(etcdEndpoint)}
            , _client{_etcdEndpoint} {
        }

        std::string const _etcdEndpoint;
        etcd::Client _client;
    };

}// namespace cfgStore
