echo off & color 0A
setlocal EnableDelayedExpansion
set DIR="./proto"
echo DIR=%DIR%

echo %GOPATH%/src
set curpath= %cd% 
set path1=/proto
set path2=/third_party/proto
echo path1=%curpath:~1,-1%%path1%
echo path2=%curpath:~1,-1%%path2%


for /R %DIR% %%f in (*.proto) do ( 
echo %%f
protoc ^
-I "%curpath:~1,-1%/proto" ^
-I "%curpath:~1,-1%/third_party/proto" ^
-I "%GOPATH%/src" ^
--gocosmos_out=plugins=interfacetype+grpc,^
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. %%f
)
pause