## Why this is valuable

This experiment is part of an overall strategy to interconnect the Foundation's longterm goals with the strategic technology decisions needed to reach them. We aim to uncover, explore, discuss and define the key challenges in "modernizing". More importantly, we aim to illuminate the opportunities.

##### To meet the Foundation's goal
[Knowledge as a service](https://meta.wikimedia.org/wiki/Strategy/Wikimedia_movement/2017)

"To serve our users, we will become a platform that serves open knowledge to the world across interfaces and communities. We will build tools for allies and partners to organize and exchange free knowledge beyond Wikimedia. Our infrastructure will enable us and others to collect and use different forms of free, trusted knowledge."

##### delivered by
[Modernizing the Wikimedia product experience](https://meta.wikimedia.org/wiki/Wikimedia_Foundation_Medium-term_plan_2019)

"We will make contributor and reader experiences useful and joyful; *moving from viewing Wikipedia as solely a website, to developing, supporting, and maintaining the Wikimedia ecosystem as a collection of knowledge, information, and insights with infinite possible product experiences and applications*."

##### the mission is to architect   
- a modern platform
- that can serve collections of knowledge and information,
- created from multiple trusted sources,
- to nearly-infinite product experiences
- and other platforms.

Success depends on:
- understanding the tradeoffs, especially in areas that have been philosophically off limits
- enabling a continuous flow of informed decisions
- cross-functional discovery and iterative step taking
- defining aspirational terms like "modern"
- understanding the cost

By "cost", we mean money, eventually. We also mean time, energy and expertise. We mean social and cultural changes that may remove roadblocks. And we mean discerning the balance between investment and goals. When they aren't balanced, how shall we alter those goals?  

### In this proof of value (PoV)
We built a prototypical first step in this direction. Our step may be too big -- or not big enough, leaving too many other questions unanswered. Either way, we are purposefully triggering the discussions that will help us discern. More questions than answers. The primarily value of this work is what we learn and how we will apply that learning.

The first step is a doozy, regardless. There is no *iterative* path to transforming a semitruck (monolith) into a small fleet of fast cars (a "modern" system). The workflow needed to build cars is quite different from the one keeping the truck on the road.

In this deliverable, we've created:
- A tiny [modern platform](#modern-platform)
- that can serve [collections of knowledge](#collections-of-knowledge)
- created from [multiple trusted sources](#multiple-trusted-sources)
- to many [product experiences](#many-product-experiences) and [other platforms](#other-platforms).

What happens when we move knowledge out of a Wikipedia website into this modern platform? In this case, we used simple wikipedia and did some limited experimentation with select English wikipedia pages. There are *many* [challenges to consider](#challenges-to-consider) before this prototype is "production ready". Production ready was not our goal. We are [engaging with some of those challenges next](#next-steps).

If you'd like to skip right to the details, read the [implementation overview](#implementation-overview) and handy-wavy [caveats](#caveats) to see our implementation choices and [how to view the demo](#demo).

## Modern platform
We've made some choices about what "modern platform" means. The wider world of content / knowledge systems are facing similar challenges ... how to "create once, publish everywhere"? We draw thinking from others, about the emerging patterns and challenges. We also relied on 18 months of [architectural explorations](#todo-add-link-with-links-and-summary) conducted prior to this PoV. These enabled us to understand what *we* need from a "modern" platform.

We define modern platform as interrelated capabilities relying on emerging industry patterns (see below). At a system level, these patterns are "implementation details" because they lay the foundation for low-level interactions between sources and products that scale as the system scales. We need this. Focusing on patterns and their implementation enables emergence, interrelated capabilities becoming greater than the sum of their parts. While the parts remain stable, predictable, changeable and encapsulated.

The patterns we've explored include:

*Canonical data modeling*: allows content to be understood by people, programs and machines outside the traditional boundaries of MediaWiki. And, as far as possible, allows consumers to request only what they need.

What is the structure of "knowledge" and how does it flow across the system? Building this data model requires defining boundaries around current parts (pages, especially) and their interrelationships. For example, a page has parts (sections) and is also part of collections (about the same topic). In our modeling, we included:

- Define a predictable structure[5] using industry-standard formats like schema.org (to support predictability and reusability)
- Break down prexisting structures (all the content on the Philadephia page) into parts (a section on the History of Philadephia) and establish interrelationships between the parts (to support "only what they need") using hypermedia linking.
- Enhance the structure with contextual information by associating parts with Wikidata (to enable natural collections like US Cities) and indexing collections with Elasticsearch.
- Enable interaction with the structure via API calls. Multiple API calls can be wrapped into a single payload -- or not.

[Working draft of our CDM](todo-edit-and-link)

*Essential note: Honestly, we don't know if it's humanly possible to "structure" Wikipedia documents. There are many questions to consider. We are identifying the Biggest Challenges so we can raise them and resolve them organizationally.

*Loose coupling*: New ways to interact with, enhance or process content (capabilities) that operate independently and are built on top of (or adjacent to) the data model.

*Event-based interactions*: activities in the system happen only when they need to happen (asynchronously) with only the information they need to accomplish their aim.

*CQRS*: Differentiating between reading and editing. In the PoV, the current structure inside of MediaWiki is left alone, it is the "trusted source". When changes happen in MW, the new system reacts by getting the necessary information and translating it into the canonical data model. This means the design works for reading but not for editing. If > 90% of the requests are for reads, can editing be a separate part of the system? We're looking at the editing workflow next.

### In the PoV, we ...
- Respond to an edit event triggered by MediaWiki when page content has changed
- Retrieve the content from the source, break it down into sections and give each part a predictable structure
- Save the page and sections as json object associated by hypermedia links (hasparts/isparts)
- Return requests for pages and/or parts (or some of page or part)
- Send objects for third-party topic analysis (the Rosette service returns a list of topics and their weight)
- Save the topics associated with the parts they describe
- Return requests for parts associated with the topic (highest-scoring sections on Physics, for example)
- Though this list seems sequential, these activities are asynchronous

## Collections of knowledge and information
Our goal was to focus on the knowledge. During our architectural explorations, a single blocker arose again and again. At the heart of our ecosystem, the knowledge we want to share with the world exists only in the context of "web page". At the moment, it is a gigantic, tangled, monolithically-orchestrated bundle of proprietary text. This text has enabled *terrific* benefit to the world. And, it's safe to say, won't meet the system's longterm goals. A predictable data model is (to some extent) needed to feed multitudes of new and varied product experiences.

Some of the knowledge is held in pages within pages (within pages) related loosely by unique software logic. Instructions for displaying knowledge on a web page are inextricably mixed into the knowledge we hope to distribute beyond a Wikipedia page. How do we structure knowledge that can shift context - from a web page to Alexa, for example? Exchanging free knowledge beyond Wikimedia requires a predictable language of exchange, beyond HTML.

The page, like a body without a skeleton, has no predictable shape until MediaWiki pieces the bones together to form a Wikipedia web page. Products and platforms that consume knowledge outside the context of MediaWiki need knowledge structured as distributable, consumable information.

As we move towards loose coupling, parts of the system will be for editing, other parts will be enhancing and consuming collections of knowledge and information. Decoupling enables parts to be built and operate independently. Yet, decoupling depends the knowledge itself being shared in a software-and-context-agnostic way. How much "structure" is enough? For example, does a page about a person have sections based on schema.org recommendations? Should "sections" or "references" exist as a structure in Mediawiki (breaking down the "body blob" as other CMS software has done)?

When we say collections, what do we mean? A page is one, predominant, collection of knowledge and information. What are the others? What is the shape of knowledge that can be consumed by "infinite product experiences"? Experiences that will likely control how the knowledge is displayed and how users interact with it.

The knowledge also needs cataloging, predictable, prebuilt relationships between knowledge that don't rely on extensive fuzzy logic queries. Like finding a book according to it's Dewey Decimal Section rather than searching the library for it. Knowledge about a subject, like Barak Obama, can be shared from wherever it currently lives on Wikipedia.

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

## Implementation overview

### Caveats

### Demo
