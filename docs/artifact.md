## To contribute

- Push changes ||
- Comment in PR ||
- Set up in StackEdit to share (instructions added the first time we do this) ||
- Paste in Google Doc and add comments

## Areas of focus
Add any notes or subjects to include here

- What
- Challenges
-- Breaking down
-- Expanded use cases (will it get used)
-- Infrastructure
-- Conceptual buy in
-- Is CQRS a viable pattern? (editing)
- Next steps
- How
-- How to use the demo
-- What is here and how it works (models)

## Foundation Goal
[Knowledge as a service](https://meta.wikimedia.org/wiki/Strategy/Wikimedia_movement/2017):

"To serve our users, we will become a platform that serves open knowledge to the world across interfaces and communities. We will build tools for allies and partners to organize and exchange free knowledge beyond Wikimedia. Our infrastructure will enable us and others to collect and use different forms of free, trusted knowledge."

Delivered through [Modernizing the Wikimedia product experience](https://meta.wikimedia.org/wiki/Wikimedia_Foundation_Medium-term_plan_2019)

"We will make contributor and reader experiences useful and joyful; *moving from viewing Wikipedia as solely a website, to developing, supporting, and maintaining the Wikimedia ecosystem as a collection of knowledge, information, and insights with infinite possible product experiences and applications*."

## Our Goal
Architect a  
- modern platform
- that can serve collections of knowledge and information,
- created from multiple trusted sources,
- to nearly-infinite product experiences
- and other platforms.

Some guiding questions are: What does that look like? How do we get from here to there? What are the challenges? Are they possible to resolve? Are they worth the necessary investment?

## This proof of value
We call this a "proof of value" (PoV) because we aren't strictly interested in proving the concept. We are also interesting in facilitating further discussion and exploration: is the cost (in time and energy) reasonably balanced against the overall goals? If not, do we alter the goals?  

This PoV is an experiment.

- Design a tiny modern platform that can serve collections of knowledge and information, created from multiple trusted sources, to many product experiences and platforms.
- Find out what happens when we move knowledge out of a Wikipedia website into this modern platform.

## Modern platform
We made some choices about what "modern platform" means. Generally speaking, we relied on emerging industry patterns for content systems facing this challenge (which is most of them). A predictable data structure, temporal decoupling The patterns include:

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
Our goal wasn't to build a rock-solid deployable prototype. A production-ready design for requires significantly more infrastructure discussion. Our goal was to focus on the knowledge. Because at the moment, it is a gigantic, tangled, monolithically-orchestrated bundle of Wikitext. Which, it's safe to say, is not going to feed the multitudes of new product experiences.

As we attempt to structure the knowledge and information for a modern platform, What are the biggest challenges we face? Are they worth the effort they'll require?

## The Experiment
During our two years of architectural explorations, a single blocker arises again and again. At the heart of our ecosystem, the knowledge we want to share with the world is a mountain of Wikitext. Like a body without a skeleton, it has no predictable shape, until MediaWiki pieces bones together to form a Wikipedia web page.

The knowledge is bounded by the context of a Wiki web page (for example, the Philadelphia page). Some of the knowledge is pages within pages related according to a unique logic. Instructions for displaying the knowledge on a Wikipedia page are inextricably mixed into the knowledge we hope to share beyond Wikimedia.

We use the word "decouple" often. Untangling the enmeshed parts is essential to modernization. In systems, the lowest level patterns scale to define the system itself. So we begin by defining some predictable structure of the knowledge as distributable information.

What is the shape of knowledge that can be consumed by "infinite product experiences"? Product experiences that will likely control how the knowledge is displayed and perhaps how users interact with it. How do we structure knowledge that can shift context - from a web page to Alexa, for example. Exchanging free knowledge beyond Wikimedia requires a predictable structure of exchange.

The knowledge also needs cataloging, a predictable relationship structure so it can be shared meaningfully. For example, as knowledge about a subject, like Barak Obama, wherever it currently lives on Wikipedia.

There are nearly-infinite challenges. Here are the primary ones we've uncovered. First, here's what the PoV does. (And the demo.)

## Change the patterns
