[![rodellison](https://circleci.com/gh/rodellison/alexa-slick-dealer.svg?style=shield)](https://app.circleci.com/pipelines/github/rodellison/alexa-slick-dealer)

Preface: My intent for forking, and ultimately making modifications to existing skill 
'Alexa Slick Dealer' is solely to use it as a means for learning how to develop an Alexa Skill using _**golang**_ and a little about
**_CircleCI_** (Continuous Integration).  **I have no intention on  releasing this skill.**  Rather, the modifications made are so that I can have a template
to use for future skills. 

What's new in this forked version?
- Type/struct enhancements in Request and Response structs for handling Display devices incorporating 
Alexa Presentation language.
- New functions replace the prior 'NewSimpleResponse' to facilitate Ask and Tell type responses'. 
- Hardcoding of _**ShouldEndSession**_ has been removed to allow for more control. 
- Functional changes in how the skill delivers data - namely, this version of the skill 
only says three items at a time, and then asks the user if they want to here more. 
(Prior the skill would just keep going through possibly dozens of deals with no stopping) 
- Session Attribute logic was added to carry forward items fetched from the frontpage/popular deals page. 
- Changed 'Simple' Card to 'Standard' Card - so as to facilitate display of images in the cards provided
back during the session. Simple cards can still be used if desired by changing the type in the response.go file
- More complete unit tests for additional coverage
- Minor MakeFile enhancements to support build/testing on OSX as well as AWS
- Moved supporting images to a subdir

**NOTE**: '**serverless.yml**' should be adjusted as appropriate if this enhanced skill is cloned or used. 



_______________________
# Slick Dealer Skill for Amazon Alexa

[Slick Dealer](https://www.amazon.com/gp/product/B07J43J36F?ie=UTF8&ref-suffix=ss_rw) is an Amazon Alexa Skill built with the Go programming language. It is an unofficial virtual assistant to the [Slickdeals](https://www.slickdeals.net) website, with no endorsements from Slickdeals.

## How it Works

While Slickdeals does not have an official API for consuming its data, it does offer a public RSS feed that contains deal titles along with descriptions and links. The RSS feed is in XML format which is easily readable with Golang.

The Slick Dealer Skill, while written in Golang, is designed to be ran as an AWS Lambda function. The responses to the function are formatted for Amazon Alexa.

## Author Information

This Skill was written by [Nic Raboy](https://www.nraboy.com), who is an advocate of modern web and mobile development technologies. He has experience in Java, JavaScript, Golang and a variety of frameworks such as Angular, NativeScript, and Apache Cordova. Nic writes about his development experiences related to making web and mobile development easier to understand.

## Resources

The Polyglot Developer - [https://www.thepolyglotdeveloper.com](https://www.thepolyglotdeveloper.com)