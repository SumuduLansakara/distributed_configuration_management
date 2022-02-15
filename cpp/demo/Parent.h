#pragma once

#include <cfgStore/LocalComponent.h>
#include <plog/Log.h>

namespace demo {
    class Parent : public cfgStore::LocalComponent {
    public:
        explicit Parent(std::string const& name)
            : cfgStore::LocalComponent("parent", name) {
        }

        void onComponentConnected(cfgStore::RemoteComponentStub stub) override {
            PLOGN << "connection handler: test: key1=" << stub.spawn().getParam("key1");
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
