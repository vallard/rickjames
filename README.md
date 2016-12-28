# Rick James [![Build Status](http://drone.ciscopipeline.io/api/badges/vallard/rickjames/status.svg)](http://drone.ciscopipeline.io/vallard/rickjames)

A Basic [sparkbot](https://developer.ciscospark.com) to deploy on Kubernetes.


Hook is: http://rickjames.glueflakes.com/spark-hook

## Initialize	
The initialize script is to deploy a webhook that starts the bot listening for events.  


Test it out: 

```
curl localhost:8080/spark-hook -X POST \
-d '{"id" : "Y2lzY29zcGFyazovL3VzL01FU1NBR0UvOWQ0OWU5NDAtY2NjMS0xMWU2LTg3YjYtZGY1OTVjOTg0NzEz" , \
"roomId" : "Y2lzY29zcGFyazovL3VzL1JPT00vYzJjNDY3MDAtYzhhMS0xMWU2LThmNmEtZTlmZTYyZjkwMzU1", \
"roomType" : "group", \
"personId" : "Y2lzY29zcGFyazovL3VzL1BFT1BMRS83MDE2MDRmMS04ZTk2LTRiNzEtOTE1Mi0wODY0YTkxYmM2MTM",\
"personEmail" : "rickjames@sparkbot.io"}' \
-H "Content-Type: application/json"
```