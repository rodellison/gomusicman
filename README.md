[![rodellison](https://circleci.com/gh/rodellison/GoMusicMan.svg?style=shield)](https://app.circleci.com/pipelines/github/rodellison/GoMusicMan)

_______________________
# The Music Man for Amazon Alexa
![Logo](Songkick108.png)

[The Music Man](https://www.amazon.com/RodEllison-The-Music-Man/dp/B01GOO5L0E) is an Amazon Alexa Skill built with the Go programming language. 
It is an unofficial virtual assistant to the [SongKick](https://www.songkick.com) website, with no endorsements from Songkick.


## How it Works

Songkick provides an API for querying Artists and Venue information. 

The Music Man Skill, written in Golang, is designed to be run as an AWS Lambda function. When invoked from an Alexa device, 
the skill parses slot content provided by Alexa (containing artist or venue information), and uses
that information to query the Songkick API for event data. Event Calendar data returned is managed within session attributes with 
each response providing 3 events at a time. The user is asked if they would like to hear more events, and if so the skill processes
another loop of 3 items. 

## Tech Notes

- To use the Songkick API, an API key is required. 
- This codebase includes provisions for Serverless deployment (serverless.yml) and CircleCI CICD  (.circleci/config.yml) 
for automated master branch testing and deployment. 
- Environmental variables are used to provide key data in both files. For the most part, with changes to 
a few Environmental variables locally, or setup directly (at CircleCI), one should be able to apply this 
code as a base for pretty much any skill.

## Author Information

This Skill was written by [Rod Ellison](https://www.rodellison.net), generated [from this template](https://github.com/rodellison/gomusicman). 

## Resources
The Music Man Alexa Skill provides information to both display oriented devices (Show, Spot, Firestick) using Alexa Presentation
Language, as well as non-display oriented devices (Echo, Dot, Tap).  
On display devices, background images
 are used from **_Songkick_**, as well as [Unsplash](https://www.unsplash.com)
 
 
Attribution to the following talented photographers on **Unsplash** for the use of their imagery:

- **Vishnu R Nair** https://unsplash.com/photos/m1WZS5ye404

- **Yvette de Wit** https://unsplash.com/photos/NYrVisodQ2M

- **Sebastian Ervi** https://unsplash.com/photos/xJt6Gs20Uqc
