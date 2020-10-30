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

### In this proof of value
We design and test a first step in the direction described above. It may be too big a step. Or not big enough, leaving too many other questions unanswered.

The challenge is: there is no *iterative* path from here to there. If the goal is to transform a semitruck into a small fleet of fast cars, while the semitruck is still in service ... best to begin with the workflow needed to build cars.

In this PoV, we:
- Designed a tiny modern platform that can serve collections of knowledge and information, created from multiple trusted sources, to many product experiences and platforms.
- Experimented with what happens when we move knowledge out of a Wikipedia website into this modern platform.

## Modern platform
In this design, we made some choices about what "modern platform" means. We draw from the wider world where content providers are facing similar challenges - how to move past the boundaries of a single piece of software towards "create once, publish everywhere"? We also relied on 18 months of architectural explorations conducted prior to this PoV to understand what *we* need from a "modern" platform.

We define modern platform as a system of interrelated capabilities relying on emerging industry patterns. The patterns of a system are what make (or break) the system as it scales. The patterns we're exploring include:

*Canonical data modeling*: allows content to be understood by people, programs and machines outside the traditional boundaries of MediaWiki. And, as far as possible, allows consumers to request only what they need.

Building this model requires defining boundaries around parts and their interrelationships. For example, a page has parts (sections) and is also part of collections (about the same topic). We included:

- Define a predictable structure[5] using industry-standard formats like schema.org (to support predictability and reusability)
- Break down prexisting structures (all the content on the Philadephia page) into parts (a section on the History of Philadephia) and establish interrelationships between the parts (to support "only what they need") using hypermedia linking.
- Enhance the structure with contextual information by associating parts with Wikidata (to enable natural collections like US Cities) and indexing collections with Elasticsearch.
- Enable interaction with the structure via API calls. Multiple API calls can be wrapped into a single payload -- or not.

LINK Working draft of our CDM

Note: Honestly, we don't know if it's humanly possible to "structure" Wikipedia documents. We are identifying the Biggest Challenges so we can raise them and resolve them organizationally.

*Loose coupling*: New ways to interact with, enhance or process content (capabilities) operate independently and are built on top of (or adjacent to) the data model.

*Event-based interactions*: activities in the system happen only when they need to happen (asynchronously) with only the information they need to accomplish their aim.

*CQRS*: The current structure inside of MediaWiki is left alone. When changes happen in MW, the new system reacts by getting the necessary information and translating it into the canonical data model. Also, if the system is used > 90% of the time for reads, can editing be a seperate part of the system?

LINK

### In the PoV, we ...
TODO: Decoupling the payload (justify the content platform design specifically)
TODO: How the implementation choices mirror these patterns.

## Collections of knowledge and information

Our goal wasn't to build a rock-solid deployable prototype. A production-ready design for requires significantly more infrastructure discussion. Our goal was to focus on the knowledge. Because at the moment, it is a gigantic, tangled, monolithically-orchestrated bundle of Wikitext. Which, it's safe to say, is not going to feed the multitudes of new product experiences.

During our two years of architectural explorations, a single blocker arises again and again. At the heart of our ecosystem, the knowledge we want to share with the world is a mountain of Wikitext. Like a body without a skeleton, it has no predictable shape, until MediaWiki pieces bones together to form a Wikipedia web page.

The knowledge is bounded by the context of a Wiki web page (for example, the Philadelphia page). Some of the knowledge is pages within pages related according to a unique logic. Instructions for displaying the knowledge on a Wikipedia page are inextricably mixed into the knowledge we hope to share beyond Wikimedia.

We use the word "decouple" often. Untangling the enmeshed parts is essential to modernization. In systems, the lowest level patterns scale to define the system itself. So we begin by defining some predictable structure of the knowledge as distributable information.

What is the shape of knowledge that can be consumed by "infinite product experiences"? Product experiences that will likely control how the knowledge is displayed and perhaps how users interact with it. How do we structure knowledge that can shift context - from a web page to Alexa, for example. Exchanging free knowledge beyond Wikimedia requires a predictable structure of exchange.

The knowledge also needs cataloging, a predictable relationship structure so it can be shared meaningfully. For example, as knowledge about a subject, like Barak Obama, wherever it currently lives on Wikipedia.

## Created from multiple trusted sources
TODO: Rosette as our example
TODO: Using schema.org (goes here?)

## Nearly-infinite product experiences

TODO: GraphQL and FE agnostic
TODO: What do those product experiences expect?

## And other platforms

TODO: Can push as well as pull

## Next steps
- Verify value (we are doing)
- Design a production-deployable version (do we have the use cases right now?)
