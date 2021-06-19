#-*- coding: utf-8 -*-
import xmind
from xmind.core.const import TOPIC_DETACHED
from xmind.core.markerref import MarkerId

w = xmind.load("test.xmind") # load an existing file or create a new workbook if nothing is found

s1=w.getPrimarySheet() # get the first sheet
s1.setTitle("first sheet") # set its title
r1=s1.getRootTopic() # get the root topic of this sheet
r1.setTitle("we don't care of this sheet") # set its title

s2=w.createSheet() # create a new sheet
s2.setTitle("second sheet")
r2=s2.getRootTopic()
r2.setTitle("root node")

# Empty topics are created from the root element and then filled.
# Examples:

# Create a topic with a link to the first sheet given by s1.getID()
t1 = r2.addSubTopic()
t1.setTopicHyperlink(s1.getID()) 
t1.setTitle("redirection to the first sheet") # set its title

# Create a topic with a hyperlink
t2 = r2.addSubTopic()
t2.setTitle("second node")
t2.setURLHyperlink("https://xmind.net") 

# Create a topic with notes
t3 = r2.addSubTopic()
t3.setTitle("third node")
t3.setPlainNotes("notes for this topic") 
t3.setTitle("topic with \n notes")

# Create a topic with a file hyperlink
t4 = r2.addSubTopic()
t4.setFileHyperlink("logo.jpeg") 
t4.setTitle("topic with a file")

# Create topic that is a subtopic of another topic
t41 = t4.addSubTopic()
t41.setTitle("a subtopic")

# create a detached topic whose (invisible) parent is the root
d1 = r2.addSubTopic(topics_type = TOPIC_DETACHED)
d1.setTitle("detached topic")
d1.setPosition(0,20)

# loop on the (attached) subTopics
topics=r2.getSubTopics()
# Demonstrate creating a marker
for topic in topics:
    topic.addMarker(MarkerId.starBlue)
    
# create a relationship    
rel=s2.createRelationship(t1.getID(),t2.getID(),"test") 

# and we save
xmind.save(w,"test2.xmind") 
