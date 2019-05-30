<p align="center">
  <img src ="https://raw.githubusercontent.com/TIBCOSoftware/flogo/master/images/flogo-ecosystem_Rules.png" />
</p>

<p align="center" >
  <b>Rules read-only tuple cache example</b>
</p>

<p align="center">
  <img src="https://travis-ci.org/TIBCOSoftware/flogo.svg"/>
  <img src="https://img.shields.io/badge/dependencies-up%20to%20date-green.svg"/>
  <img src="https://img.shields.io/badge/license-BSD%20style-blue.svg"/>
  <a href="https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link"><img src="https://badges.gitter.im/Join%20Chat.svg"/></a>
</p>

## Steps to build and run cache example
Install redis, run build and start the flogo app. Run 'cat data.txt | redis-cli --pipe' to load house tuples to redis cache. Then run 
'curl "localhost:7766/solar/eligible?parcel=0002&bill=300"' to fire the solar eligibility action. Running 
curl "localhost:7766/solar/eligible?parcel=0001&bill=300" will not fire the action as parcel 0001 already has the solar installed.
