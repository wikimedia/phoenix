## Why this is valuable

This experiment is part of an overall strategy to interconnect the Foundation's longterm technology goals with the strategic decisions needed to reach them. We aim to uncover, explore, discuss and define the key challenges in "modernizing". More importantly, we aim to illuminate opportunities.

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

By "cost", we mean money, eventually. We also mean time, energy and expertise. We mean social and cultural changes that may remove roadblocks, when necessary. And we mean discerning whether or not the effort required is reasonably balanced against the value of our goals. When it isn't, how shall we alter those goals?  

### In this proof of value (PoV)
We built a prototypical first step in the direction described above. This may be too big a step. Or not big enough, leaving too many other questions unanswered. The primarily value of this work is what we learn and how we will apply that learning.

The challenge is: there is no *iterative* path to transforming a semitruck (monolith) into a small fleet of fast cars (component parts). We are designing the workflow needed to build cars while keeping the truck on the road.

The deliverable is:
- A tiny [modern platform](#modern-platform)
- that can serve [collections of knowledge](#collections-of-knowledge)
- created from [multiple trusted sources](#multiple-trusted-sources)
- to many [product experiences](#many-product-experiences) and [other platforms](#other-platforms).

What happens when we move knowledge out of a Wikipedia website (in this case simple wikipedia and some limited experimentation with select English wikipedia pages) into this modern platform? There are *many* [challenges to consider](#challenges-to-consider) before this prototype is "production ready". We are [engaging with some of those challenges next](#next-steps).

To skip right to the details, read the [implementation overview](#implementation-overview), handy-wavy [caveats](#caveats) in our implementation choices and [how to view the demo](#demo).

## Modern platform
We made some choices about what "modern platform" means. We draw from the wider world where content systems are facing similar challenges - how to move past the boundaries of a single piece of software towards "create once, publish everywhere"? We also relied on 18 months of [architectural explorations](#todo-add-link-with-links-and-summary) conducted prior to this PoV that enabled us to understand what *we* need from a "modern" platform.

We define modern platform as interrelated capabilities relying on emerging industry patterns (see below). These patterns lay the foundation for low-level interactions between sources and products to scale as the system scales. The goal is to enable emergence, interrelated capabilities becoming greater than the sum of their parts, while the parts remain stable, predictable, changeable and encapsulated.

The patterns we're exploring include:

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

## Collections of knowledge and information

Our goal wasn't to build a rock-solid deployable prototype. A production-ready design requires *significantly* more infrastructure discussion (which we are also having). Our goal was to focus on the knowledge. Because at the moment, it is a gigantic, tangled, monolithically-orchestrated bundle of Wikitext. This has enabled terrific benefit. And, it's safe to say, can't feed multitudes of new and varied product experiences without predictable data modeling.

During our pre-PoV architectural explorations, a single blocker arose again and again. At the heart of our ecosystem, the knowledge we want to share with the world exists only in the context of "web page". The knowledge is bounded by the context of a Wiki web page (for example, the Philadelphia page). Some of the knowledge is pages within pages related according to a unique logic. Instructions for displaying the knowledge on a Wikipedia page are inextricably mixed into the knowledge we hope to share beyond Wikimedia.

A page is one, predominant, collection of knowledge and information. What are the others? What is the shape of knowledge that can be consumed by "infinite product experiences"? Product experiences that will likely control how the knowledge is displayed and perhaps how users interact with it. How do we structure knowledge that can shift context - from a web page to Alexa, for example. Exchanging free knowledge beyond Wikimedia requires a predictable structure of exchange.

The page, like a body without a skeleton, has no predictable shape. MediaWiki pieces the bones together to form a Wikipedia web page. To serve products and platforms that consume knowledge outside the context of MediaWiki, we need to predictably structure the knowledge as distributable information. Designing this will be a back and forth dialogue over time: is this structure used when content is created? (For example, do we recommend sections based on the type of article?) Or added after it is saved (without pushing parsing onto consumers)?

As we move towards loose coupling, parts of the system will be editing, enhancing and consuming collections of knowledge and information. Decoupling these parts enables them to be built and operate independently. Decoupling depends the knowledge itself being shared in a software-and-context-agnostic way.  

The knowledge also needs cataloging, predictable relationship structures, so it can be shared meaningfully. For example, sharing knowledge about a subject, like Barak Obama, wherever it currently lives on Wikipedia.

## Multiple trusted sources
The PoV uses simplewikipedia as the primary source but the same pattern would apply to any source. As long as we can respond to an event by getting the change. [Rosette](TODO link) is the source for topics (what is this object about?) and other sources can be added. The topics are Wikidata items, which can also be a source for further information about the topic.

## Many product experiences

TODO: GraphQL and FE agnostic
TODO: What do those product experiences expect?

## And other platforms

TODO: Can push as well as pull

## Challenges to consider
TODO

## Next steps
- Verify value (we are doing)
- Design a production-deployable version (do we have the use cases right now?)

## Implementation overview

### Caveats

### Demo
