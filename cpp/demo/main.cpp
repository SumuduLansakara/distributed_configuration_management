#include "Child.h"
#include "Parent.h"
#include <cfgStore/Communicator.h>

#include <plog/Appenders/ColorConsoleAppender.h>
#include <plog/Formatters/TxtFormatter.h>
#include <plog/Initializers/RollingFileInitializer.h>
#include <plog/Log.h>
#include <regex>
#include <unistd.h>

using namespace std;

int main() {
    static plog::ColorConsoleAppender<plog::TxtFormatter> consoleAppender;
    plog::init(plog::debug, &consoleAppender);

    cfgStore::Communicator::init(cfgStore::ETCD_CONN_STR);
    demo::Parent c1{"comp1"};
    c1.initConfigs({{"key1", "value1"}});
    c1.connect();
    c1.watchComponents("child");

    demo::Child c2{"comp2"};
    c2.initConfigs({{"key1", "value1"}});
    c2.connect();
    c2.setParam("newParam", "newVal2");
    c2.setParam("newParam", "newVal3");
    c2.removeParam("newParam");

    //
    std::cout << " >-- " << std::endl;
    for (auto& kv : demo::Parent::listComponents("parent")) {
        std::cout << kv.first << ":" << std::endl;
        for (auto& v : kv.second)
            std::cout << " - " << v << std::endl;
    }
    std::cout << " <-- " << std::endl;
    //

    c2.disconnect();
    c1.disconnect();

    PLOGD << "exiting";
}
