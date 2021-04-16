del a*
go build -o a.dll -buildmode=c-shared .
gcc -o main c\main.c a.dll
.\main.exe