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

## Goal
Knowledge as a service is essential to the [2030 strategy](https://meta.wikimedia.org/wiki/Strategy/Wikimedia_movement/2017):

"To serve our users, we will become a platform that serves open knowledge to the world across interfaces and communities. We will build tools for allies and partners to organize and exchange free knowledge beyond Wikimedia. Our infrastructure will enable us and others to collect and use different forms of free, trusted knowledge."

Modernizing the Wikimedia product experience is one of two goals in the [Wikimedia Foundation Medium-term plan 2019](https://meta.wikimedia.org/wiki/Wikimedia_Foundation_Medium-term_plan_2019). Modernizing our product experience:

"We will make contributor and reader experiences useful and joyful; moving from viewing Wikipedia as solely a website, to developing, supporting, and maintaining the Wikimedia ecosystem as a collection of knowledge, information, and insights with infinite possible product experiences and applications."

This experiment has one foundational question: if we move knowledge out of the Wikipedia website and into a modern platform ... what happens?

We are breaking down the HTML page into pre-structured knowledge that can be shared with people, machines and products. This immediately solves some challenges. It also creates many new ones.

Our goal with this Proof of Value is to facilitate the next phase of exploration: What are those challenges and are they worth tackling?

## The Experiment
During our two years of architectural explorations, a single blocker arises again and again. At the heart of our ecosystem, the knowledge we want to share with the world is a mountain of Wikitext. Like a body without a skeleton, it has no predictable shape, until MediaWiki pieces bones together to form a Wikipedia web page.

The knowledge is bounded by the context of a Wiki web page (for example, the Philadelphia page). Some of the knowledge is pages within pages related according to a unique logic. Instructions for displaying the knowledge on a Wikipedia page are inextricably mixed into the knowledge we hope to share beyond Wikimedia.

We use the word "decouple" often. Untangling the enmeshed parts is essential to modernization. In systems, the lowest level patterns scale to define the system itself. So we begin by defining some predictable structure of the knowledge as distributable information.

What is the shape of knowledge that can be consumed by "infinite product experiences"? Product experiences that will likely control how the knowledge is displayed and perhaps how users interact with it. How do we structure knowledge that can shift context - from a web page to Alexa, for example. Exchanging free knowledge beyond Wikimedia requires a predictable structure of exchange.

The knowledge also needs cataloging, a predictable relationship structure so it can be shared meaningfully. For example, as knowledge about a subject, like Barak Obama, wherever it currently lives on Wikipedia.

There are nearly-infinite challenges. Here are the primary ones we've uncovered. First, here's what the PoV does. (And the demo.)

## Change the patterns
