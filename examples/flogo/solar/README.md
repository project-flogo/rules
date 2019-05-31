<p align="center">
  <img src ="https://raw.githubusercontent.com/TIBCOSoftware/flogo/master/images/flogo-ecosystem_Rules.png" />
</p>

## Solar eligibility check


This flogo app check if a house a eligible for solar panel marketing promotion. house data that contains parcel number and the status of solar panel installation. solar check event contains a parcel number and monthly spending bill. 

The app loads the house data from redis cache storage and, as solar check events arrives, fires the solar eligibility action depending the condition described as follows:

If the monthly bill is greater than 200 and the house doesn't have a solar panels installed,
the house with the matching parcel id should be a candidate for solar panel installation promotion.

house tuples are to be pre-loaded from data.txt 
The command in data.txt loads tuples with a tuple key as a redis hash key. It then updates the index named after
the tuple type by associating the hash key to the index.

<house data tuples>
house:parcel:0001 parcel 0001 is_solar true
house:parcel:0002 parcel 0002 is_solar false

solar events are asserted for a monthly electiricity bill generated.

<solar event tuples>
solar:parcel:0001 parcel 0001 bill 300
solar:parcel:0002 parcel 0002 bill 250


## Steps to build and run example
Install redis, build and start the flogo app. <br/>
Run <br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;'cat data.txt | redis-cli --pipe' to load house tuples to redis cache. <br/>

Then run <br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;curl "localhost:7766/solar/eligible?parcel=0002&bill=300" <br/>
to fire the solar eligibility action. <br/>

Running <br/> 
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;curl "localhost:7766/solar/eligible?parcel=0001&bill=300" <br/>
will not fire the action as parcel 0001 already has the solar installed.
