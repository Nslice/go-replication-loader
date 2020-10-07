go build -ldflags "-s -w"

go-replication-loader -prjName "Testing Client" ^
-c CMS_TestingClient ^
-d "C:\Program Files (x86)\eLeed_TestingClient" ^
-u admin ^
-p 12 ^
-f sergey.zalunin@akforta.com ^
-smtp mail.akforta.com ^
-smtplogin sergey.zalunin ^
-smtppass vgy78uhb ^
-dbname TestingClient ^
-backuppath "F:\AutoBackup" ^
-usecompr ^
-saveargs ^
-rsd ^
-t sergey.zalunin@akforta.com ^
-t olesya.mezger@akforta.com ^
-t ritika@screen.ae ^
-t kallidus@akforta.com ^
-t dhivakar@screen.ae ^
-t surbhi@screen.ae ^
-t talgat.tairov@akforta.com

go clean