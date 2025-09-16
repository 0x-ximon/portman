cmake_minimum_required(VERSION 4.0)
project(linux)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)


find_package(Qt5 COMPONENTS
        Core
        Gui
        Widgets
        REQUIRED)

add_executable(linux 
        main.cpp)
        
target_link_libraries(linux
        Qt5::Core
        Qt5::Gui
        Qt5::Widgets
)

