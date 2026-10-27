# Tell Me About Yourself

The self-introduction that opens almost every interview, built with the
**Present → Past → Pivot** structure.

## How to use this file

- **It's an index, not a story.** The intro points to my proof (the STAR
  stories in `star-stories.md`); it doesn't tell them. Whatever I mention
  is what the interviewer will ask about next, so I choose the agenda.
- **Memorize the order, not the text.** Speak from the skeleton below.
  Only the lines marked ★ and the short transition phrases are memorized
  word for word: the first line (nerves are highest at the start) and the
  last line (the only sentence I fully control).
- **Every technology I name is a promise.** Only name what I can defend
  for two minutes without notes. RabbitMQ, idempotency, retries and DLQ
  are in; gRPC and OpenTelemetry stay out until I've practiced them.
- **Every sentence must survive "tell me more".** No inflated numbers
  ("twelve clients", not "dozens"), no claims narrower or wider than the
  truth ("across the whole ERP", not "mostly financial").
- **Never end with "...and yeah, that's it".** The last sentence hands the
  floor back with a hook about the company.
- **Shorten by cutting blocks, never by speaking faster.**

## Core version (60–100 s, ~220 words)

> **Present.** I'm a backend engineer at Teraware, a Brazilian ERP
> company, where I've worked fully remote for five years. Over that time
> I've delivered features end to end for twelve clients, across the whole
> ERP, from CRM to financials, plus integrations with outside systems.
> The project I'm proudest of is a bank reconciliation system built on
> Open Finance APIs, which is still in production today.
>
> **Past.** My background is in Computer Engineering. I tutored C for a
> year, so pointers and memory were never new to me. In parallel, I spent
> two years at the Apple Developer Academy, shipping apps to the App Store
> in teams of five, with designers. That's where I learned to think about
> the product, not just the code.
>
> **Pivot.** Over the last few months, I've been focusing on Go, because
> it fits what I already build: services that talk to other services and
> have to be reliable when money is involved. To go deep, I built an
> order execution engine in Go. A REST API publishes events to RabbitMQ,
> and a consumer processes them idempotently, with retries, and a
> dead-letter queue for messages that keep failing. It's deployed and
> running. Now I'm looking for a role where Go is the main language. Your
> team stood out to me because you're building payment infrastructure,
> which is the kind of problem I already solve with bank reconciliation.

## Skeleton (the only thing I look at while speaking)

```
PRESENT
  ★ "I'm a backend engineer at Teraware, a Brazilian ERP company,
     where I've worked fully remote for five years."
  end to end · TWELVE clients · whole ERP · CRM → financials · integrations
  proudest → bank reconciliation · Open Finance · still in prod

PAST  ("My background is in Computer Engineering.")
  tutored C → pointers + memory never new
  in parallel → Apple Academy · 2 yrs · App Store · teams of 5 + designers
  "That's where I learned..." → product, not just code

PIVOT  ("Over the last few months, I've been focusing on Go,")
  why → fits what I BUILD (services ↔ services · money)
  order engine → REST → RabbitMQ → consumer · idempotent
               → retries, AND a DLQ for messages that keep failing · deployed
  ★ "Now I'm looking for a role where Go is the main language."
  ★ "Your team stood out to me because [why them]."
```

## Short version (~30 s, ~70 words)

Same skeleton with cuts: the whole Past block goes, the breadth line in
Present goes, and the engine shrinks to keywords.

> I'm a backend engineer with five years of fully remote experience at a
> Brazilian ERP company, building features end to end for twelve clients,
> including a bank reconciliation system on Open Finance APIs that's still
> in production. Over the last few months I've been focusing on Go: I
> built an event-driven order execution engine with RabbitMQ, idempotent
> processing and a dead-letter queue, and it's deployed. Now I'm looking
> for a role where Go is the main language.

If a "tell me more" follows, don't restart: continue with the Past block
("Sure. My background is in Computer Engineering...").

## Which version, when

| Situation | Version |
|---|---|
| Start of a coding or system design round | Short (it's a warm-up; extra time comes out of the problem) |
| Recruiter asks for "a quick intro" | Short |
| Outreach message or Gophers Slack intro | Short, written |
| Recruiter screen, hiring manager, behavioral round | Core |
| English/communication step on vetting platforms, founder call | Core + one ownership line (sole engineer for some clients) |

## The "why them" line

The last sentence is the only one that changes per company, and it's
written **before** the interview, never improvised. Each application in
the tracking sheet gets a "Why them" column with this sentence ready.

- **Fintech / payments:** "...because you're building payment
  infrastructure, which is the kind of problem I already solve with bank
  reconciliation."
- **Sports tech:** "...because live match data is a real-time,
  event-driven problem, and that's exactly what I've been building toward.
  And I'm a huge football fan, so that helps."

## Likely follow-ups

**"Five years — was all of that full-time engineering?"**
The first two years were as an intern, but I was already writing code
that went to production. For the last three, I've been full-time, and I'm
the only engineer responsible for some of our clients.

**"Why Go? Why move away from PHP?"**
It's not that I'm running away from PHP; it taught me a lot about
production systems. It's that the problems I solve most, integrations and
financial data processing, are exactly where Go shines: built-in
concurrency, static typing, and a single binary that's easy to deploy. I
want to work with the right tool for the problems I already enjoy
solving.

**"How long have you been working with Go?"**
My first contact was a side project, a football betting pool website with
real users. Since July, I've been studying Go in a structured way, every
day, and the order execution engine came out of that. In production at
work, my stack is still PHP and PostgreSQL.

**"Tell me about your education."**
I graduated in Computer Engineering from a federal university whose
program ranks among the top three in Brazil, with an average above 9 out
of 10. I tutored C for a year and was part of the programming club,
solving algorithm problems every week.

**"What does your production stack look like?"**
PHP and PostgreSQL in production, and about 70% of my work is backend;
the rest is the frontend for my own features. Go is where I'm building
now.

## Kept out on purpose

| Left out | Why | Where it goes |
|---|---|---|
| Personal life (age, relationship, moving) | Not for the first minute | Small talk, if it comes up |
| Any negative reason for leaving | Badmouthing the current job is the most cited mistake | Never, in that form |
| MySQL → PostgreSQL migration and the 40% | It was an assigned task, and the number isn't defended yet | A STAR story, after verifying the number |
| Client names | Confidentiality | STAR stories, as "a football club", "a law firm" |
| gRPC, OpenTelemetry, Swift | Naming is promising | Technical follow-ups |
| University ranking and GPA | Written beats spoken: the listener can't weigh them | Resume (with the source) |
| The university acronym | Means nothing to a foreign interviewer | Resume |

## Pronunciation

- **ERP:** letter by letter, "ee-ar-pee".
- **reconciliation:** "rek-un-sil-ee-AY-shun".
- **tutored:** "TOO-terd".
- **Academy:** stress on "KA", "uh-KA-duh-mee".
- **RabbitMQ:** "RA-bit em-KYOO".
- **dead-letter queue:** "ded-LE-ter KYOO".
- **Initial R** (*role, RabbitMQ, REST, reliable, retries,
  reconciliation*): the Portuguese initial "r" sounds like an English "h",
  so *role* can be heard as *hole*. English "r": the tongue tip rises and
  curls back slightly without touching anything. Quick test: dictate "I'm
  looking for a role" to phone voice typing; if it writes "hole", practice
  the "r".

## Practice log

| Date | Take | Time | Notes |
|---|---|---|---|
| 2026-10-26 | 1 | 1:52 | Order of blocks and transitions right; Present in 28 s (on target). Issues: "dozens of clients", "retries in a dead-letter queue" (wrong order), "Go fits what I already use" (should be *build*), "why Go" missing, no closing line. Long pauses in the Past block (two cause-and-effect links in one sentence). 4 fillers. |
| 2026-10-26 | 2 | 2:10 | All structural fixes landed: Past flowed, "why Go" present, DLQ order correct, closing line delivered. Regressions: "fully remote" dropped, "dozens" again. C mentioned twice; closing sentence too long (~35 words) → pauses at the end. 4 fillers. Fixes: "fully remote" moved into the ★ opening line, "a dozen" → "twelve", closing split into two sentences, C only in Past. |
