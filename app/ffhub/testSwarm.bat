cls
@echo off
cd c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\ffhub\changestatus\
echo reinstall changestatus
go install

cd c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\ffhub\inputtracker\
echo reinstall inputtracker
go install

cd c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\ffhub\manager\
echo reinstall manager
go install

cd c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\ffhub\profiler\
echo reinstall profiler
go install

cd c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\ffhub\fflocator\
echo reinstall fflocator
go install

cd c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\ffhub\
echo done

for /f %%i in ('gum choose --header="Chto zapuskaem?" --header.foreground="150" inputtracker manager changestatus profiler fflocator') do set TEST_PROG=%%i
echo run %TEST_PROG%:


%TEST_PROG%
