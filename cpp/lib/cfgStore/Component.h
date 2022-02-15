#pragma once

#include "Constants.h"
#include <etcd/Client.hpp>
#include <etcd/Watcher.hpp>
#include <fmt/core.h>
#include <utility>

namespace cfgStore {

    class Component {
    public:
        Component(std::string kind, std::string name,
                  std::string instanceId)
            : _kind{std::move(kind)}
            , _name{std::move(name)}
            , _instanceId{std::move(instanceId)} {}

        virtual ~Component() = default;

        std::string const& getKind() { return _kind; }
        std::string const& getInstanceId() { return _instanceId; }
        std::string const& getName() { return _name; }
        std::string prettyName() { return fmt::format("{} [{}]", _name, _instanceId); }

        virtual void setParam(std::string const& key, std::string const& val) = 0;
        virtual std::string getParam(std::string const& key) = 0;
        virtual void removeParam(std::string const& key) = 0;

    protected:
        std::string const _kind;
        std::string const _name;
        std::string const _instanceId;
    };

}// namespace cfgStore
