# Architectural Proof of Value
### by the Architecture Team

[Mural board of this artifact](https://app.mural.co/t/neweditorexperiences1471/m/neweditorexperiences1471/1606310676597/5994555cfce1cc40ff0b0ec85600da818c7bea12)

## Why this artifact is valuable

The architecture team interconnects the Foundation's strategically-imperative goals with the technology decisions needed to reach them. This artifact describes a prototyping exercise demonstrating the value of modern architectural patterns for knowledge systems. We also uncover and describe potentially-disruptive challenges blocking the modernization of our knowledge system.

##### To meet the Foundation's goal
"To serve our users, we will become a platform that serves open knowledge to the world across interfaces and communities. We will build tools for allies and partners to organize and exchange free knowledge beyond Wikimedia. Our infrastructure will enable us and others to collect and use different forms of free, trusted knowledge." -- [Knowledge as a service](https://meta.wikimedia.org/wiki/Strategy/Wikimedia_movement/2017)

##### delivered by
"We will make contributor and reader experiences useful and joyful; *moving from viewing Wikipedia as solely a website, to developing, supporting, and maintaining the Wikimedia ecosystem as a collection of knowledge, information, and insights with infinite possible product experiences and applications*." -- [Modernizing the Wikimedia product experience](https://meta.wikimedia.org/wiki/Wikimedia_Foundation_Medium-term_plan_2019)

##### our mission is to architect   
- a modern platform
- that can serve collections of knowledge and information,
- created from multiple trusted sources,
- to nearly-infinite product experiences
- and other platforms.

### In this proof of value (PoV)
We created a prototypical step towards a modern platform:

- A tiny, experimental [modern platform](#modern-platform)
- that can serve [collections of knowledge](#collections-of-knowledge)
- created from [multiple trusted sources](#multiple-trusted-sources)
- to many [product experiences](#many-product-experiences) and [other platforms](#other-platforms).

#### The four leverage points
However we approach it, the first step is a doozy. There is no *iterative* path towards transformation. Neither is there a lift-and-shift migration option. We need to find leverage points: capabilities in the system that we can decouple from the current day-to-day operations. As challenging as leverage points may be, to find and to change, they unlock highly-valuable opportunities. While simultaneously laying a strong and cohesive foundation for the future system.

The leverage points explored in this PoV are:
- giving shape and structure to Knowledge
- designing inherent relationships between knowledge parts to create collections
- building decoupled relationships between parts of the system rather than building capabilities into the software (this includes changing the choreography of essential activities)

#### The big questions
The PoV leaves many questions unanswered -- on purpose. We are *triggering cross-functional discussions and decisions needed* to discern a path forward. While we have more questions than answers, we are significantly more confident in the questions. Top four include:

1. What is "just enough" structure needed for the knowledge?
2. What infrastructure can support these patterns at scale?
3. From a system point of view, can reading be decoupled from editing?
4. How will modernization impact the current editing workflow?

The highest-value next step is continuing to gather and apply learnings from teams across the foundation that help answer these questions.

There are *many* [challenges to consider](#challenges-to-consider) before this prototype is "production ready". *Production ready was not our goal.* We are [engaging with some of those challenges next](#next-steps).

If you'd like to skip right to the details, read the [implementation overview](#implementation-overview) and handy-wavy [caveats](#caveats) to see our implementation choices and [how to view the demo](#demo).

## Modern platform
We've made choices about what "modern platform" means. These choices were informed by the wider world of content and knowledge systems, which are facing similar challenges. How to "create once, publish everywhere"? How do we *distribute* knowledge to wherever people engage with it?

We drew thinking from others about emerging patterns and challenges. We also relied on 18 months of [architectural explorations](#todo-add-link-with-links-and-summary) conducted prior to this exercise. These explorations enabled us to identify what *we* need from a "modern" platform. Some needs are in synch with the world at large and a few (essential challenges) are unique.

We define modern platform as interrelated capabilities relying on emerging industry patterns (see below). In the PoV and at a system level, these patterns are the implementation details. They lay the foundation for low-level interactions between knowledge sources and products that scale as the system scales.

### Patterns are the thing
Patterns enable us to design for emergence: create interrelated capabilities that become greater than the sum of their parts. So we focused on patterns that enable parts that are stable, predictable, changeable and encapsulated. Patterns that let us design a system by focusing on ...
- the data model of "knowledge"
- the parts that deliver the capabilities (things the system does)
- the relationship between the parts
- and the structure of their interaction

The patterns we've explored include:

*Canonical data modeling*: allows content/knowledge to be understood by people, programs and machines outside the traditional boundaries of MediaWiki. And, as far as possible, allows consumers to request only what they need.

What is the structure of "knowledge" and how does it flow across the system? Building this data model requires defining boundaries around data objects and their interrelationship. A page, for example, is a collection of sections. (And templates, which we did not tackle here.) Sections are also part of collections about a topic (physics, for example.) In our modeling, we:

- Defined a predictable structure[5] using industry-standard formats like schema.org (to support predictability and reusability)
- Broke down prexisting structures (all the content on the Philadephia page) into parts (a section on the History of Philadephia) and establish interrelationships between the parts (to support "only what they need") using hypermedia linking.
- Enhanced the structure with contextual information by associating parts with Wikidata (to enable natural collections like US Cities) and indexing collections with Elasticsearch.
- Enabled interaction with the structure via API calls. Multiple API calls can be wrapped into a single payload -- or not.

[Working draft of our CDM](todo-edit-and-link)

Essential note: Honestly, we don't know if it's humanly possible to "structure" Wikipedia content sufficiently. There are *many* questions to consider. We identify the Biggest Challenges so we can raise them and resolve them organizationally.

*Loose coupling*: New ways to interact with, enhance or process content (capabilities) that operate independently and are built on top of (or adjacent to) the data model.

*Event-based interactions*: activities in the system happen only when they need to happen (asynchronously) with only the information they need to accomplish their aim.

*CQRS*: Differentiating between reading and editing. In the PoV, the current structure inside of MediaWiki is left alone, it is the "trusted source". When changes happen in MW, the new system reacts by getting the necessary information and translating it into the canonical data model. This means the design works for reading but not for editing. If > 90% of the requests are for reads, can editing be a separate part of the system? We're looking at the editing workflow next.

## Implementation overview
- Respond to an edit event triggered by MediaWiki when page content has changed
- Retrieve the content from the source, break it down into sections and give each part a predictable structure
- Save the page and sections as json object associated by hypermedia links (hasparts/isparts)
- Return requests for pages and/or parts (or some of page or part)
- Send objects for third-party topic analysis (the Rosette service returns a list of topics and their weight)
- Save the topics associated with the parts they describe
- Return requests for parts associated with the topic (highest-scoring sections on Physics, for example)
- Though this list seems sequential, these activities are asynchronous
- We used simple wikipedia and did some limited experimentation with select English wikipedia pages.

TODO: Add model

### Demo

For the month of January, 2020, you can access a demo instance [here]
(TODO). The demo is:

- a front end that interacts with GraphQL and the structured content store.
- the structured content store contains content from Simple Wikipedia, updated when edits are made there.
- topics associated with each object (page, section) from [Rosette](TODO add link).

### Caveats

## Collections of knowledge and information
How can we design knowledge to be consumed by "infinite product experiences"? How do we enable these "experiences" to control how the knowledge is displayed and how users interact with it. When we say collections, what do we mean? A page is one, predominant, type of collection of knowledge. What are the others?

### The shape of knowledge
During our architectural explorations, a single blocker arose again and again. At the heart of our ecosystem, the knowledge we want to share with the world isn't made for modern distributions. It exists as a "web page" made from a gigantic, tangled, monolithically-orchestrated bundle of proprietary text.

This bundle of text has enabled *terrific* benefit to the world. The challenge is, without detangling, it won't meet the system's longterm goals in the emerging digital world.

A predictable data model is needed (to some extent) to feed multitudes of new and varied product experiences. Products and platforms that consume knowledge outside the context of MediaWiki need the knowledge structured as distributable, consumable information.

How much "structure" is enough? For example, should about a page about a person have sections based on schema.org recommendations? This would make it more consumable by products and platforms. Should "sections" or "references" exist as a structure of knowledge, inside and outside MediaWiki?

At the moment, some of the knowledge is in pages within pages (within pages), related loosely by unique software logic. The page, like a body without a skeleton, has no predictable shape until MediaWiki pieces the bones together to form a Wikipedia web page. Instructions for displaying knowledge on a web page are inextricably woven into the knowledge we hope to distribute beyond a Wikipedia page. How do we enable knowledge to shift context - from a web page to Alexa, for example?

Exchanging free knowledge beyond Wikimedia requires loose coupling and a predictable language of exchange, beyond HTML. Loose coupling enables parts of the system to be built and operate independently. Editing, for example, doesn't need to be enmeshed with reading. Machine learning can provide necessary information without being ensconced inside editing software. Decoupling depends on the knowledge itself being shared in a software-and-context-agnostic way.

### The shape of collections
A page is a container for a collection of parts. A page was the initial shape, a website was a framework that related containers to other containers, thus building a collection of knowledge. Now, we are building a new framework that includes pages but is not defined by them.

Categories are collections. Pages and parts of pages associated with a Wikidata item is a collection. Collections are relationships developed, programatically or by editors, between pieces of knowledge. The way humans envision and plan these relationships shapes the way the knowledge is developed.

To form scalable collections, the knowledge needs cataloging. Consistency of relationships between knowledge parts makes collections consumable by nearly-infinite products and platforms. Without overtaxing the system with queries. Predictable, prebuilt relationships that don't rely on extensive fuzzy logic queries are ideal.

## Multiple trusted sources
The PoV uses simplewikipedia as the primary knowledge source. But the same pattern will apply to adding any subsequent source. There can be multiple wikipedia's, for example. The platform responds to an event sent from the source by getting the change from the source's API. As long as both are possible, a source is likely a valid option.

[Rosette](TODO link) is our source for topics, creating collections based on what the knowledge is about. Other context-creating sources can be added similarly. The topics from Rosette are Wikidata items. Wikidata can also be a source to enhance information about the topic.

## Many product experiences and other platforms
When we imagine "nearly-infinite product experiences", what comes to mind? Answering that question is cross-functional work happening now. For the PoV, we imagined things like:

- product experiences requesting knowledge so they can build their own "page", or collection or context for displaying.
- these experiences drawing from multiple sources and needing relationships that give the integrated knowledge meaning (about Barack Obama, for example)
- a website or app about Cricket that draws people towards Wikipedia in places that aren't part of the community yet
- any decoupled frontend experience

For platforms, we imagined
- big platforms who use the free knowledge getting exactly what they need (and perhaps monetizing that request)
- interrelationships with What's App and Facebook that draw people into learning and perhaps editing
- pushing knowledge to platforms

For a deeper dive into the product experiences and platforms that are supported by modernization, see [PoV products and platforms](TODO).

## Challenges to consider
The primary challenge is embracing uncertainty. We can't know exactly how an emergent system will emerge. We can make well-reasoned decisions while exploring the path to modernization. We can design for emergence. But this requires embracing uncertainty and making sound decisions in the midst of it.

Modeling and planning for change triggers confusion and anxiety, two things that will most certainly push the system in the wrong direction (to regain status quo). This is a challenge must not be underestimated.

Other major challenges include:

- Ensuring the infrastructure will scale
- Agreeing on "just enough structure" of the current page content
- Breaking down a page into a data model that differs from the crowdsourced version
- Versioning (across the system)
- Understanding how creating knowledge in the "sources" interrelates with serving that knowledge everywhere

## Next steps
The next steps are delivering four further artifacts:

- Infrastructure exploration -- designing a production-viable version of this AWS prototype
- Models of current editing workflows, so we understand them
- Editing exploration - what must change in the knowledge sources in order to decouple editing from many products and platforms? Are they viable?
- Creating the architecture repository as a space to understand and explore emerging systems patterns

Teams are also creating front end prototypes demonstrating uses for these patterns. Details on which teams are not yet decided.

Many branches of discussion have already begun. Their success depends on:
- understanding the tradeoffs, especially in areas that have been philosophically off limits
- enabling a continuous flow of informed decisions
- cross-functional discovery and iterative step taking
- defining aspirational terms like "modern"
- understanding the cost

By "cost", we mean estimating the financial investment, though its too soon for that analysis. We also mean articulating the time, energy and expertise required. When we say cost, we also mean the social and cultural changes that may be necessary to remove roadblocks. And discerning the balance between our values, goals and investments.
