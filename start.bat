go build

.\go-replication-loader.exe -prjName "Testing Client" ^
-c "CMS_Project" ^
-d "H:\eLeed\eLeed" ^
-u admin ^
-dbname Eleed_Head_Origin ^
-backuppath "C:\drive" ^
-dbuserid sa ^
-dbpassword WTF_MR_JOE$75 ^

go clean