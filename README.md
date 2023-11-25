ical-git is a lightweight calendar application written in Go. 

The calendar data is a collection of iCal files under version control. ical-git is a deamon that periodically pulls the iCal files and creates notifications based on the VALARM components within the iCal files. ical-git currently supports notifications through Telegram bots and local Linux desktop notifications.


