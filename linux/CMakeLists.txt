cmake_minimum_required(VERSION 4.0)
project(portman_linux)

find_package(Qt5 COMPONENTS 
    Core 
    Gui 
    Widgets 
    REQUIRED
)

add_executable(portman_linux src/main.cpp)
target_include_directories(portman_linux PUBLIC include lib)
        
target_link_libraries(portman_linux
    Qt5::Core
    Qt5::Gui
    Qt5::Widgets
)

