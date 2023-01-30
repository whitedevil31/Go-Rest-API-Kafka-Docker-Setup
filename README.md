# Problem statement 

The product is a form creation portal where users can create different forms and specify the input types and attach 
post submission triggers which needs to be executed after the form is submitted in a fail-safe manner and should scale to million request . 
The logic should be designed in such a way thgat the upcoming triggers should be easily  plugged into the existing code and consistency should be maintained 

# Preparing the solution 
  ## TASK 1
  
  Design an efficient way to model the questions and responses of the forms and
  how I would store them in the data store

  ## TASK 2
  
   Efficiently execute the post-submission triggers (the events that the customer
   expects to run after a form is submitted)

   ## TASK 3
   
   Design a logic such that new upcoming triggers should be easily be plugged into
   the existing code and consistency should be maintained
   
   
## POSTMAN COLLECTION AND TECH STACK 
    Postman collection of the API - https://www.postman.com/collections/c22d05a8c75dba16d75f
    ## Tech stack followed 
    Golang - Just trying out ...
    MongoDB - MongoDB is my go-to DB for any database interaction
    Kafka -To handle event triggers
    Docker to containerize the application


# SOLVING TASK #1
So initially when creating a form we will add all the questions which are to be required in the
form.The schema of the form is described below

   ![image](https://user-images.githubusercontent.com/63295747/215503239-4e63a666-63f3-4619-8591-6c0e7c5992da.png) 
    
Once a user adds all the required questions and fills necessary information like QuestionTitle
,AnswerType which could be optional/required depending upon the use-case and clicks on
the create form button in the frontend the form will be created.
The backend logic will be simple. We will have a REST API Endpoint of POST method (
http://localhost:8080/api/add_form ) to handle the request and save the formData .We will
store the form data in the “forms” collection.
When a form is opened its data can be fetched from the forms collections using the formId
.Once it is filled and submit button is clicked ,we will have a REST API Endpoint of POST
method(http://localhost:8080/api/submit_form ) to handle the submission .
The schema for the submission is described below

![image](https://user-images.githubusercontent.com/63295747/215503319-e80d7c6b-3b84-4d47-8023-df1888f20694.png)


The intuition behind saving the Id of all the entities instead of actual data of the entity itself is
once we have the ID we can use the query/aggregation feature of MongoDB to bring all the
required additional data .We can even compute analytics for a particular answer/questions
through this _id. The _id of MongoDB enables faster query over other keys and also by this
we don’t consume redundant data onto the disk .We will save this submissionData in the
“submissions” collection.
For handling the list of events to trigger for a specific form we will maintain an “events”
collection where we save events along with the formId .
We can add events by sending a POST request to (http://localhost:8080/api/add_event)
and get all events for a form by sending a GET request to
(http://localhost:8080/api/get_events/{formId} )

![image](https://user-images.githubusercontent.com/63295747/215503387-3a46847c-5d0a-403f-a8da-e833d8110a86.png)



# SOLVING TASK #2

Here I would like to demonstrate the scenario by executing these two events after our form is
submitted .
EVENT 1 : Add the submission result in a Google sheet

EVENT 2: Send a Email notification to the user to confirm his form submission (I picked
Email service instead of SMS as there aren’t enough free SMS service providers and both
Email and SMS examples are almost similar )

Initially I thought of simply executing these events in the form of conventional functions or
HTTP CALLS (Assuming these event logics are written in an external API ) then given the
condition that the solution needs to be failsafe,scalable,fault tolerance lead me to a
conclusion that this can’t be a ideal solution . It took me a day to find out about Apache
Kafka and I believe it could provide a really good solution here .
Kafka acts as our event streamer here . Once our form response is stored in the
database,our primary API will act as a producer-api and it will send messages to the
consumer-api and return the response .

The consumer will execute the functions(trigger
events) then in its own server . By this we don’t need to wait for the events to finish to return
the HTTP response from our producer.
I found multiple ways to achieve the above task

  #WAY1 : Trigger events as TOPICS | single consumer in single consumer group for all topics

  #WAY2 : Trigger events as TOPICS | single consumer in multiple consumer groups with different topics

  #WAY 3 : Trigger events as TOPICS | one consumer group with multiple consumers

  #WAY 4 :  Have a single topic | multiple consumer groups with multiple consumers with specific topics

I chose WAY 4 after researching a lot on various forums . The reason for using this approach
is that even though will be streaming on a single topic since we have multiple consumers
(number of consumers will be almost equal to number of partitions of our primary topic)
in the consumer group the traffic will be distributed and if incase we are to increase the
number of topics we can go ahead with multiple consumers in multiple consumer groups
with different topics .

One more reason why I didn’t go ahead with the other approaches was the lack of
dynamic topic creation support in the consumers and creating a consumer for each topic
could lead to a lot of unwanted deployments in the cluster in future
By this #WAY4 approach we can expect a really good throughput with low latency .
In my codebase I have a single Kafka broker with a single producer with 2 partitions for
the topic . A Zookeeper to track our cluster information like Kafka brokers ,leader election
and configuration data of topics and partitions . Two consumers(for demonstration purpose)
in a single consumer group listening to specific topics .We can easily increase the
consumers in our consumer group by simple duplicating the main deployment of our
consumer

I have also implemented the concurrency approach of Go by creating a separate
Goroutine for a util function call and used channels to transfer data

![image](https://user-images.githubusercontent.com/63295747/215504579-1339a52c-f4de-42c3-af86-ed17067c8b79.png)


# SOLVING TASK #3

In our current flow the producer API sends messages to a single topic and these message
consists of a key,value pair where the key is the type of event to trigger and value is the
data
So based on the key I determine what type of event to trigger and call it .


![image](https://user-images.githubusercontent.com/63295747/215504806-0a425e0e-8ee3-418e-bba7-5f3c0dfaa699.png)

I used the built-in reflect package of Go to call functions based on an input string in a
dynamic fashion . Using this approach we maintain consistency of the codebase and avoid
the need of multiple if else/switch statements to alter between different functions
which could be lengthy as well as involve lot of redundant code
All the functions we need to add can be included in the events.go file with the event name as
the function name .
So the entire flow for adding a event to an form is :
STEP 1: Send a POST request to (http://localhost:8080/api/add_event) to attach an event
to a form

STEP 2 : Define the business logic in the events.go file as a function with function name as
the event name which we used in the the POST method above(refer to postman collection)
That’s it ,just two steps and that’s all we have to do to add a new trigger to our form.
By this approach, we can easily integrate all our new logic which makes the process more
generic and can be considered as just plug and play fashion

# IMPLEMENT LOGGING


There are various approaches to implement logging for our API .I went ahead to implement a
simple custom log handler which could add logs to a common file logs.txt . We can pass the
required data to log into and can even classify logs based on the context . But in production
grade applications it is advised to use a popular package since it could provide a lot of built
in features and I found Zap to be a really good fit

![image](https://user-images.githubusercontent.com/63295747/215505541-d7501540-3d4c-4ba6-b6b3-3882194cc92e.png)




