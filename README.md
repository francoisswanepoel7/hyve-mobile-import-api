Hyve Mobile Import API
--
This API is located at: https://import.tsekcorona.co.za/

Data structures:
emailValidationMap : Map of validated emails with CSV Data value struct
domainIpMap: Map of Domains to Resolved IPs

Calls:

POST:

https://import.tsekcorona.co.za/process
 
1. Coroutines: 
    - Email Validation
    - IP From Domain
    - Inserting validated data into database 
    

GET:

https://import.tsekcorona.co.za/import 
 - This imports records into the database
 - Would be run via cron 

https://import.tsekcorona.co.za/import
 - This gets processed records and posts them to https://api.tsekcorona.co.za/contact with a limit 0f 1000 records per post
 - Processed records get 'exported' flag set to true in database
 - Would be run via cron


