#pragma once

#include <cfgStore/LocalComponent.h>
#include <plog/Log.h>

namespace demo {

    class Child : public cfgStore::LocalComponent {
    public:
        explicit Child(std::string const& name)
            : cfgStore::LocalComponent("child", name) {
        }

        void onComponentConnected(cfgStore::RemoteComponentStub stub) override {
            PLOGN << "connection handler: " << stub.spawn().getParam("key1");
            watchParameters(std::move(stub));
        }

        void onComponentDisconnected(cfgStore::RemoteComponentStub stub) override {
            PLOGN << "disconnection handler: " << stub.getInstanceId();
        }

        void onParameterSet(cfgStore::RemoteComponentStub stub, std::string const& key) override {
            PLOGN << "parameter set handler: " << stub.getInstanceId() << ":" << key;
        }

        void onParameterRemoved(cfgStore::RemoteComponentStub stub, std::string const& key) override {
            PLOGN << "parameter deletion handler: " << stub.getInstanceId() << ":" << key;
        }
    };
}// namespace demo
