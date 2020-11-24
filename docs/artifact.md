# Architectural proof of value
### by the Architecture team

## Why this artifact is valuable

The architecture team aims to interconnect the Foundation's longterm goals with the strategic, technology decisions needed to reach them. To succeed, we need to uncover, explore, discuss and define the key challenges in "modernizing". More importantly, we hope to illuminate the opportunities.

This artifact describes our first step towards that goal.

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
During this PoV, we built a prototypical step towards a modern platform. Our step may be too big, or not big enough, leaving too many questions unanswered. Either way, we are purposefully triggering the discussions needed to help us discern. We end with more questions than answers, yet we are significantly more confident in the questions.

The primarily value of this work is continuing to gather and apply the learnings from teams across the foundation.

The first step is a doozy. There is no *iterative* path towards transforming a semitruck (monolith) into a small fleet of ubers (a "modern" system). The workflow needed to build an interconnected system of cars is quite different from the one keeping the truck on the road.

In this experiment, we've created:
- A tiny, experimental [modern platform](#modern-platform)
- that can serve [collections of knowledge](#collections-of-knowledge)
- created from [multiple trusted sources](#multiple-trusted-sources)
- to many [product experiences](#many-product-experiences) and [other platforms](#other-platforms).

We used simple wikipedia and did some limited experimentation with select English wikipedia pages. There are *many* [challenges to consider](#challenges-to-consider) before this prototype is "production ready". *Production ready was not our goal.* We are [engaging with some of those challenges next](#next-steps).

If you'd like to skip right to the details, read the [implementation overview](#implementation-overview) and handy-wavy [caveats](#caveats) to see our implementation choices and [how to view the demo](#demo).

## Modern platform
We've made choices about what "modern platform" means. These choices were informed by the wider world of content and knowledge systems, which are facing similar challenges. How to "create once, publish everywhere"? How do we *distribute* knowledge to wherever people engage with it?

We drew thinking from others about emerging patterns and challenges. We also relied on 18 months of [architectural explorations](#todo-add-link-with-links-and-summary) conducted prior to this exercise. These explorations enabled us to identify what *we* need from a "modern" platform. Some needs are in synch with the world at large and a few (essential challenges) are unique.

We define modern platform as interrelated capabilities relying on emerging industry patterns (see below). In the PoV and at a system level, these patterns are the implementation details. They lay the foundation for low-level interactions between knowledge sources and products that scale as the system scales.

### Patterns are the thing
We need to focus on patterns and how to implement them. They enable us to design for emergence: create interrelated capabilities that become greater than the sum of their parts. While the parts remain stable, predictable, changeable and encapsulated. Fundamentally, we mean create the
- parts: capabilities, things they system does)
- the relationship between the parts (when do they need each other?)
- and the structure of their interaction (knowledge as an universal language).

The patterns we've explored include:

*Canonical data modeling*: allows content/knowledge to be understood by people, programs and machines outside the traditional boundaries of MediaWiki. And, as far as possible, allows consumers to request only what they need.

What is the structure of "knowledge" and how does it flow across the system? Building this data model requires defining boundaries around data objects and their interrelationship. A page, for example, is a collection of sections. (And templates, which we did not tackle here.) Sections are also part of collections about a topic (physics, for example.) In our modeling, we:

- Defined a predictable structure[5] using industry-standard formats like schema.org (to support predictability and reusability)
- Broke down prexisting structures (all the content on the Philadephia page) into parts (a section on the History of Philadephia) and establish interrelationships between the parts (to support "only what they need") using hypermedia linking.
- Enhanced the structure with contextual information by associating parts with Wikidata (to enable natural collections like US Cities) and indexing collections with Elasticsearch.
- Enabled interaction with the structure via API calls. Multiple API calls can be wrapped into a single payload -- or not.

[Working draft of our CDM](todo-edit-and-link)

Essential note: Honestly, we don't know if it's humanly possible to "structure" Wikipedia content sufficiently. There are *many* questions to consider. We want to identify the Biggest Challenges so we can raise them and resolve them organizationally.

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

TODO: Add model

### Demo

### Caveats

## Collections of knowledge and information
What is the shape of knowledge that can be consumed by "infinite product experiences"? Experiences that will likely control how the knowledge is displayed and how users interact with it. When we say collections, what do we mean? A page is one, predominant, collection of knowledge and information. What are the others?

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
The PoV uses simplewikipedia as the primary knowledge / content source but the same pattern applies to any source we add. Multiple wikis, for example. As long as the platform can respond to an event sent from the source by getting the change from the source's API, it can join.

[Rosette](TODO link) is our source for topics (what is this knowledge about?). Other sources can be added that give context. The topics from Rosette are Wikidata items. Wikidata can also be a source to enhance information about the topic.

## Many product experiences
When we imagine "nearly-infinite product experiences", what comes to mind? Answering that question is cross functional work happening now. For the PoV, we imagined things like:
- product experiences requesting knowledge so they can build their own "page", or collection or context for displaying.
- these experiences drawing from multiple sources and needing relationships that give the integrated knowledge meaning (about Barack Obama, for example)
- a website or app about Cricket that draws people towards Wikipedia in places that aren't part of the community yet

## And other platforms
We imagined
- big platforms who use the free knowledge getting exactly what they need (and perhaps monetizing that request)
- interrelationships with What's App and Facebook that draw people into learning and perhaps editing
- pushing knowledge to platforms

## Challenges to consider
The primary challenge as we take these steps is embracing uncertainty. We can't know how it'll all come together. We can make well-reasoned decisions about how to explore the path. Other major challenges are:

- Ensuring the infrastructure will scale
- Agreeing on "just enough structure" of the current page content
- Breaking down a page into a data model that differs from the crowdsourced version
- Versioning (across the system)

## Next steps
The next steps are four further artifacts:

- Exploration of infrastructure (design a production version)
- Model the editing workflow (current, so we understand the needs)
- Editing exploration - is it possible to decouple editing from serve many of the "read" requests?
- Continue creating the architecture repository as a space to understand and explore emerging systems patterns

Success depends on:
- understanding the tradeoffs, especially in areas that have been philosophically off limits
- enabling a continuous flow of informed decisions
- cross-functional discovery and iterative step taking
- defining aspirational terms like "modern"
- understanding the cost

By "cost", we mean estimating the financial investment, though its too soon for that analysis. We also mean articulating the time, energy and expertise required. By cost, we also mean the social and cultural changes that may be necessary to remove roadblocks. And we mean discerning ways to balance our investment, values and goals.
