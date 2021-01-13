Architectural Proof of Value
by the Architecture Team

[Mural board of this artifact](https://app.mural.co/t/neweditorexperiences1471/m/neweditorexperiences1471/1606310676597/5994555cfce1cc40ff0b0ec85600da818c7bea12)

<!-- TOC START min:1 max:3 link:true asterisk:false update:true -->
- [Why this artifact is valuable](#why-this-artifact-is-valuable)
- [Our mission](#our-mission)
  - [To meet the Foundation's goal](#to-meet-the-foundations-goal)
  - [... we will](#-we-will)
  - [... by architecting](#-by-architecting)
  - [... beginning with a proof of value (PoV)](#-beginning-with-a-proof-of-value-pov)
- [The outcome](#the-outcome)
  - [Overview](#overview)
  - [Implementation](#implementation)
  - [Caveats](#caveats)
  - [Demo](#demo)
    - [Demo: Fetch a part of a page](#demo-fetch-a-part-of-a-page)
    - [Demo: Fetch sections by topic](#demo-fetch-sections-by-topic)
    - [Demo: GraphQL sandbox](#demo-graphql-sandbox)
- [Architecting the mission](#architecting-the-mission)
  - [Modern platform](#modern-platform)
    - [Patterns](#patterns)
    - [Leverage points](#leverage-points)
    - [Big questions](#big-questions)
  - [Collections of knowledge and information](#collections-of-knowledge-and-information)
    - [The shape of knowledge](#the-shape-of-knowledge)
    - [The shape of collections](#the-shape-of-collections)
  - [Multiple trusted sources](#multiple-trusted-sources)
  - [Many product experiences and other platforms](#many-product-experiences-and-other-platforms)
- [Challenges to consider](#challenges-to-consider)
- [Next steps](#next-steps)
<!-- TOC END -->

# Why this artifact is valuable

Architecture interconnects the Foundation's strategically-imperative goals with the system-level decisions needed to reach them. This artifact, and the prototype exercise it describes, interconnects the Foundation's strategically-imperative goal, Knowledge as a Service, with the technology decisions needed to reach it.

This artifact also demonstrates the value of modern systems patterns, uncovers leverage points (places to make impactful changes) and describes potentially-disruptive challenges blocking mission-critical decisions.

# Our mission

## To meet the Foundation's goal
"to serve our users, we will become a platform that serves open knowledge to the world across interfaces and communities. We will build tools for allies and partners to organize and exchange free knowledge beyond Wikimedia. Our infrastructure will enable us and others to collect and use different forms of free, trusted knowledge." -- [Knowledge as a service](https://meta.wikimedia.org/wiki/Strategy/Wikimedia_movement/2017)

## ... we will
"make contributor and reader experiences useful and joyful; *moving from viewing Wikipedia as solely a website, to developing, supporting, and maintaining the Wikimedia ecosystem as a collection of knowledge, information, and insights with infinite possible product experiences and applications*." -- [Modernizing the Wikimedia product experience](https://meta.wikimedia.org/wiki/Wikimedia_Foundation_Medium-term_plan_2019)

## ... by architecting
- a modern platform
- that can serve collections of knowledge and information,
- created from multiple trusted sources,
- to nearly-infinite product experiences
- and other platforms.

## ... beginning with a proof of value (PoV)
We created a prototypical step towards a modern platform.

# The outcome
## Overview
Here is [a hands-on working prototype](https://wikimedia.github.io/phoenix/) and the [implementation behind it](https://github.com/wikimedia/phoenix). The goal was to create:

- A tiny, experimental [modern platform](#modern-platform)
- that can serve [collections of knowledge](#collections-of-knowledge)
- created from [multiple trusted sources](#multiple-trusted-sources)
- to many [product experiences and other platforms](#many-product-experiences-and-other-platforms).

The demo has been shared across the organization. This small experiment has evolved into deeper explorations across the organization, including the Core Platform Team, Product Strategy Working Group, Okapi and SDAW.

If you'd like to skip right to the details, read the [implementation overview](#implementation) and handy-wavy [caveats](#caveats) or [how to view the demo](#demo).

While we created something tangible, our essential focus is on systems analysis. We are laying the foundation for systems architecture -- a practice that will support the work ahead. This work includes [designing system patterns](#patterns), discovering [leverage points](#leverage-points) and identifying the [Big Questions](#big-questions).

## Implementation
The prototype is built in AWS using SNS messages, Lambdas written in Go, S3, GraphQL, DynamoDB and Elastic Search. [See component list](https://docs.google.com/document/d/13ycCf8Mhxfs9K-hHXvsD2j1J8GcO4FQhKu4qCyx7fOs). It interacts with [Rosette](https://www.rosette.com/capability/topic-extractor/) to analyze the sections and return topics. The [repository is in github](https://github.com/wikimedia/phoenix). 

Though this list seems sequential, these activities are asynchronous.

### Event-driven workflow

Respond to a change event:
- Respond to an event stream message sent by Simple Wikipedia when article has changed
- Retrieve the article via the Parsoid API and save the raw result
- Break the raw result down into sections associated by hypermedia links - a page has parts (sections), a section is part of a page. ('hasparts/isparts')
- Save them as individual json objects with [a predictable structure](https://docs.google.com/spreadsheets/d/1ZWuczQQ0XpzCYS92PKXpIP3FM4ds0XPQyz7q9xR5GuE/edit#gid=0) using schema.org
- Save the page title associated with the resource ID

When a new file is saved:
- Send section to Rosette via API
- Save the list of resulting (most-salient) Wikidata Items (topics) associated with that section

### Request-driven workflow

- Return requests for pages and/or sections with only the data requested
- Return requests for sections associated with a topic 

[Initial diagram](https://app.lucidchart.com/lucidchart/f283e649-cdb6-4275-9452-7114571a82e7/view?page=Q3nNnx6PpfFM#) and TODO: Add Eric's updated diagram

## Caveats
There are *many* [challenges to consider](#challenges-to-consider) before this prototype is "production ready". *Production ready was not our goal.* We are [engaging with some of those challenges next](#next-steps).

## Demo
*For the month of February, 2020, you can access a demo instance [here](https://wikimedia.github.io/phoenix/)*. The demo is a front end that interacts with GraphQL and the structured content store. The structured content store contains content from Simple Wikipedia, updated when edits are made there. It also includes the topics associated with each object (page, section) from [Rosette](https://www.rosette.com/).

The demo provides several examples of potential behavior of the PoV. You can:
- fetch a section of an article by it's name
- fetch the sections associated with a specific topic
- input a GraphQL query

### Demo: Fetch a part of a page
You can fetch a section by it's name.

This example showcases how an article that was divided into sections allows for flexibly fetching the section names, and individual sections only. When a page is chosen from the drop-down list, a GraphQL query is sent to the content store, requesting the list of section titles that are available for the article. The query sets up a request for the article name, modification date, and the name of all its parts (section titles) from the requested article and populates the second drop down.

After the request is completed, the second drop down is populated, allowing the demo user to request a specific section. The second query sets up a request for the specific part inside the article itself. This means that the payload that is sent includes only the requested part, without receiving the complete article, and without expecting the consumer to process or manipulate the received content to present what is needed.

#### GraphQL queries
The queries used in this part of the demo show how easy it is to request and receive only the specific pieces of information that the consumer requires, and reduces the load of processing or manipulating the page by the consumer.

**Requesting a list of sections within the article:**

```
{
  page(name: { authority: "simple.wikipedia.org", name: "PAGE NAME"} ) {
    name
    dateModified
    hasPart(offset: 0) {
      name
    }
  }
}
```

**Requesting a specific section by name:**

```
{
  node(name: { authority: "simple.wikipedia.org", pageName: "PAGENAME”, name: "SECTION NAME" } ) {
    dateModified
    name
    unsafe
  }
}
```

### Demo: Fetch sections by topic
(TODO: This demo isn’t implemented yet, because our topic fetching isn’t implemented.)

You can request sections that are associated with a specific topic (Wikidata item).

This example shows the connection between parts (article sections) and semantic topics (wikidata items) produced by Rosette. The demo collects Rosette topics that are associated with sections and provides them in a drop-down list. Choosing a topic results in producing a GraphQL query that requests the sections that are associated with that topic, and presents them to the user. Each section then showcases the most salient topics of itself, allowing the user to explore content by wikidata topics.

#### GraphQL query

The query used to fetch sections by a given topic

Xxxxxxxx


The query used to fetch a specific section its top 5 most relevant topics

Xxxxx

### Demo: GraphQL sandbox
You can input a GraphQL query of your own.

Finally, the demo includes a GraphQL sandbox for testing and exploring the way queries are built and the payload that they produce. The sandbox is based on [GraphiQL](https://github.com/graphql/graphiql/tree/main/packages/graphiql#readme).

On the left side of the screen, the GraphQL editor allows the user to insert a custom query of their choosing, based on the available types. The user can learn what types they can request using the “Docs” popup on the top right corner of the screen.

The “Docs” button opens a popup for the “Documentation explorer”, allowing the user to view the available GraphQL hierarchical entity definition, and use those to request information.

Clicking on the “Execute query” button at the top left of the screen will run the query, and the resulting payload is presented at the right side of the interface.

The user can also look at the history of previous queries and adjust those to experiment with the different payload results.

# Architecting the mission

## Modern platform
- *A [modern platform](#modern-platform)*
- that can serve [collections of knowledge](#collections-of-knowledge)
- created from [multiple trusted sources](#multiple-trusted-sources)
- to many [product experiences and other platforms](#many-product-experiences-and-other-platforms).
We've made choices about what "modern platform" means. These choices were informed by the wider world of content and knowledge systems, where others face similar challenges. We explored emerging patterns and challenges. How do you "create once, publish everywhere"? How do we *distribute* knowledge to wherever people engage with it?

We also relied on 18 months of [architectural explorations](#TODO-add-link-with-links-and-summary) conducted prior to this exercise. These explorations enabled us to identify what *we* need from a "modern" platform. Some needs are in synch with the world at large and a few (essential challenges) are unique.

We define modern platform as interrelated capabilities relying on emerging industry patterns (see below). At a system level, these patterns are the implementation details. They lay the foundation for low-level interactions between knowledge sources and products that scale as the system scales. Which is essential to our mission.

### Patterns
Patterns enable us to design for emergence: create interrelated capabilities that can become greater than the sum of their parts. We focused on patterns that enable stable, predictable, changeable and encapsulated parts. Patterns that let us design a system by focusing on:
- the data model (the shape of) "knowledge"
- the parts that deliver the necessary capabilities (things the system does)
- the relationship between those parts
- and the structure of their interaction

The patterns we've explored include:

*Canonical data modeling*: allows content/knowledge to be understood by people, programs and machines outside the traditional boundaries of MediaWiki. And, as far as possible, allows consumers to request only what they need.

What is the structure of "knowledge" and how does it flow across the system? Building this data model requires defining boundaries around data objects and their interrelationship. A page, for example, is a collection of sections. (And templates, which we did not tackle here.) Sections are also part of collections about a topic (physics, for example.) In our modeling, we:

- Defined a predictable structure[5] using industry-standard formats like schema.org (to support predictability and reusability)
- Broke down prexisting structures (all the content on the Philadephia page) into parts (a section on the History of Philadephia) and establish interrelationships between the parts (to support "only what they need") using hypermedia linking.
- Enhanced the structure with contextual information by associating parts with Wikidata (to enable natural collections like US Cities) and indexing collections with Elasticsearch.
- Enabled interaction with the structure via API calls. Multiple API calls can be wrapped into a single payload -- or not.

[Working draft of our CDM](todo-edit-and-link)

*Loose coupling*: New ways to interact with, enhance or process content (capabilities) that operate independently and are built on top of (or adjacent to) the data model.

*Event-based interactions*: activities in the system happen only when they need to happen (asynchronously) with only the information they need to accomplish their aim.

*CQRS*: Differentiating between reading and editing. In the PoV, the current structure inside of MediaWiki is left alone, it is the "trusted source". When changes happen in MW, the new system reacts by getting the necessary information and translating it into the canonical data model. This means the design works for reading but not for editing. If > 90% of the requests are for reads, can editing be a separate part of the system? We're looking at the editing workflow next.

### Leverage points
The scope of modernization -- transforming the the world's largest reference website into the world's largest knowledge system -- is monumental. To understand where to focus our time and attention, we've identified three leverage points.

> "Folks who do systems analysis have a great belief in “leverage points.” These are places within a complex system where a small shift in one thing can produce big changes in everything." -- Donella Meadows

However we approach it, the first step is a doozy. There is no *iterative* path towards transformation. Neither is there a lift-and-shift migration option. We need to find capabilities in the system that we can decouple from the current day-to-day operations. As challenging as leverage points may be to find and to change, they unlock highly-valuable opportunities. While simultaneously laying a strong and cohesive foundation for the future system.

The leverage points explored in this PoV are:
1.  [*Giving shape and structure to Knowledge*](#the-shape-of-knowledge): Honestly, we don't know if it's humanly possible to "structure" Wikipedia content sufficiently. the knowledge we want to share with the world isn't made for modern distributions. We must try. Also, knowledge is currently shaped by the context of "web page" and that doesn't fit emerging contexts.
2.  [*Designing inherent relationships between knowledge parts to create collections*](#the-shape-of-collections): Collections are relationships developed, programatically or by editors, between pieces of knowledge. The way humans envision and plan these relationships shapes the way the knowledge is developed. The PoV pre-builds the knowledge payload (an answer to the queries) based on the relationships we know are the most valued. How would we expand this over time? 
3.  *Building decoupled relationships between parts of the system* rather than building capabilities into the software: This includes changing the choreography of essential activities ... in many ways, the paradigm itself is changing.

Exploring patterns and identifying leverage points helped us prioritize questions to explore next.

### Big questions
The scope of questions we need to answer, some we have not yet discovered, is equally monumental. The PoV leaves many questions unanswered -- on purpose. We are *triggering cross-functional discussions and decisions needed* to discern a path forward. While we have more questions than answers, we are significantly more confident in the questions. Top four include:

1. What is "just enough" structure needed for the knowledge?
2. What infrastructure can support these patterns at scale?
3. From a system point of view, can reading be decoupled from editing?
4. How will modernization impact the current editing workflow?

The highest-value next step is continuing to gather and apply learnings from teams across the foundation that help answer these questions.

## Collections of knowledge and information
- A [modern platform](#modern-platform)
- *that can serve [collections of knowledge](#collections-of-knowledge)*
- created from [multiple trusted sources](#multiple-trusted-sources)
- to many [product experiences and other platforms](#many-product-experiences-and-other-platforms).

How can we design knowledge to be consumed by "infinite product experiences"? How do we enable these "experiences" to control how the knowledge is displayed and how users interact with it? When we say collections, what do we mean? 

A page is one, predominant, type of collection of knowledge. What are the others?

### The shape of knowledge
During our architectural explorations, a single blocker arose again and again. At the heart of our ecosystem, the knowledge we want to share with the world isn't made for modern distributions. It exists as a "web page" made from a gigantic, tangled, monolithically-orchestrated bundle of proprietary text.

This bundle of text has enabled *terrific* benefit to the world. The challenge is, without detangling, it won't meet the system's long-term goals in the emerging digital world.

A predictable data model is needed (to some extent) to feed multitudes of new and varied product experiences. Products and platforms that consume knowledge outside the context of MediaWiki need the knowledge structured as distributable, consumable information.

How much "structure" is enough? For example, should about a page about a person have sections based on schema.org recommendations? This would make it more consumable by products and platforms. Should "sections" or "references" exist as a structure of knowledge, inside and outside MediaWiki?

At the moment, some of the knowledge is in pages within pages (within pages), related loosely by unique software logic. The page, like a body without a skeleton, has no predictable shape until MediaWiki pieces the bones together to form a Wikipedia web page. Instructions for displaying knowledge on a web page are inextricably woven into the knowledge we hope to distribute beyond a Wikipedia page. How do we enable knowledge to shift context - from a web page to Alexa, for example?

Exchanging free knowledge beyond Wikimedia requires loose coupling and a predictable language of exchange, beyond HTML. Loose coupling enables parts of the system to be built and operate independently. Editing, for example, doesn't need to be enmeshed with reading. Machine learning can provide necessary information without being ensconced inside editing software. Decoupling depends on the knowledge itself being shared in a software-and-context-agnostic way.

### The shape of collections
A page is a container for a collection of parts. A page was the initial shape, a website was a framework that related containers to other containers, thus building a collection of knowledge. Now, we are building a new framework that includes pages but is not defined by them.

Categories are collections. Pages and parts of pages associated with a Wikidata item is a collection. Collections are relationships developed, programatically or by editors, between pieces of knowledge. The way humans envision and plan these relationships shapes the way the knowledge is developed.

To form scalable collections, the knowledge needs cataloging. Consistency of relationships between knowledge parts makes collections consumable by nearly-infinite products and platforms. Without overtaxing the system with queries. Predictable, prebuilt relationships that don't rely on extensive fuzzy logic queries are ideal.

## Multiple trusted sources
- A [modern platform](#modern-platform)
- that can serve [collections of knowledge](#collections-of-knowledge)
- *created from [multiple trusted sources](#multiple-trusted-sources)*
- to many [product experiences and other platforms](#many-product-experiences-and-other-platforms).

The PoV uses Simple Wikipedia as the primary knowledge source. But the same pattern will apply when adding any subsequent source. There can be multiple Wikipedia's, for example. The platform would add a service that responds to an event sent from the source and gets the change from the source's API. As long as both are possible, a source is likely a valid participant.

[Rosette](https://www.rosette.com/capability/topic-extractor/) is our source for topics, creating collections based on related Wikidata items. Other context-creating sources can be added similarly. Wikidata can also be a source to enhance information about the topic.

## Many product experiences and other platforms
- A [modern platform](#modern-platform)
- that can serve [collections of knowledge](#collections-of-knowledge)
- created from [multiple trusted sources](#multiple-trusted-sources)
- *to many [product experiences and other platforms](#many-product-experiences-and-other-platforms)*.

When we imagine "nearly-infinite product experiences", what comes to mind? Answering that question is cross-functional work happening now. For the PoV, we imagined things like:

- product experiences requesting knowledge so they can build their own "page", or collection or context for displaying.
- these experiences drawing from multiple sources and needing relationships between them that give the knowledge meaning (everything about Barack Obama, for example)
- a website or app about Cricket that draws people towards Wikipedia in places that aren't part of the community yet
- any decoupled frontend experience

For platforms, we imagined
- big platforms who use the free knowledge getting exactly what they need (and perhaps monetizing that request)
- interrelationships with What's App and Facebook that draw people into learning and perhaps editing
- pushing knowledge to platforms

We also imagined the Internet of Things, News and Information sources and products designed to increase engagement.

# Challenges to consider
The primary challenge is embracing uncertainty. We can't know how this emergent system will emerge. We can make sound, well-reasoned decisions while exploring the path to modernization. We can *design* for emergence with architectural best practices. But uncertainty is our companion and we will be making sound decisions in the midst of it.

Modeling and planning for change triggers confusion and anxiety, two things that will most certainly push the system in the wrong direction (to regain status quo). This is a challenge must not be underestimated.

Other major challenges include:

- Ensuring the infrastructure will scale
- Agreeing on "just enough structure" of the current page content
- Breaking down a page into a data model that differs from the crowdsourced version
- Versioning (across the system)
- Understanding how creating knowledge in the "sources" interrelates with serving that knowledge everywhere

# Next steps
The next steps are delivering four further artifacts:

- Infrastructure: designing a production-viable version of this AWS prototype.
- Models of current editing workflows, so we understand them.
- Exploring what changes are needed in editing in order to decouple editing from many products and platforms. Is this even possible?
- Educating ourselves on modern system approaches.
- Creating the architecture repository as a space to understand and explore emerging systems patterns.

Teams like Structured Data and Okapi are eploring adopting this work as part of their future plans. 

Many branches of discussion have already begun. Their success depends on:
- understanding the tradeoffs, especially in areas that have been philosophically off limits
- enabling a continuous flow of informed decisions
- cross-functional discovery and iterative step taking
- defining aspirational terms like "modern"
- understanding the cost

By "cost", we mean estimating the financial investment, though its too soon for that analysis. We also mean the time, energy and expertise required. We mean the social and cultural changes that may be necessary to remove roadblocks. And we mean discerning the balance between our values, goals and investments.
