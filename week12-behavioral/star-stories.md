# STAR Stories

Behavioral interview prep, written in the STAR format
(Situation → Task → Action → Result).

## How to use this file

- **Memorize the skeleton, not the text.** In an interview, tell the
  story from the bullet points in your own words. A recited answer
  sounds robotic.
- **Target length: 1:45–2:10 spoken** (~260–300 words).
- **Use "I" for my actions.** Credit the team when it applies, but make
  clear which part was mine.
- **Every number must survive "how do you know?"** If I can't say where a
  number came from, it doesn't go in the story.
- **Details are follow-ups, not part of the story.** Keep the story at the
  level of decisions; each detail left out is a follow-up where I can
  go deep.

| # | Story | Main competencies | Status |
|---|---|---|---|
| 1 | Bank reconciliation (ERP) | Technical challenge, ownership, trade-offs | Done |
| 2 | Teammate left out of decisions (Apple Developer Academy) | Conflict, leadership without authority, self-awareness | Done |
| 3 | The report that took the system down (ERP) | Accountability, learning from failure, scaling | Done |
| 4 | Tight deadline | Prioritization | To do |
| 5 | Disagreeing with someone on the team | Communication, influence | To do |
| 6 | Full ownership of a problem | Ownership | To do |

---

## Story #1: Bank reconciliation

**Answers questions like:**
- "What's the biggest technical challenge you've solved?"
- "Tell me about a complex problem you solved."
- "Tell me about a project you're proud of."
- "Tell me about a technical trade-off you made."

### Skeleton (this is what I look at)

```
S (~45 words)
  - ERP, ~5 clients; no reconciliation → users compared screen by screen
  - soccer club canteen: match days → spike of small transactions → 1h+

T (~25 words)
  - I designed and built bank reconciliation end to end, on my own, ~2 months

A (~150 words)
  1. Ingest: many sources (OFX, PDF, Open Finance via Pluggy)
     → one internal format; duplicates blocked by the DB, per bank account
  2. Rank top 5 open bills or invoices (date + amount) → user confirms
     WHY: silent wrong match >> few seconds of review
     only predictable cases (fees) automated via user rules
  3. Hardest part: many-to-many, WITH A LINK TABLE
     → warning when amounts don't add up
  (+1 sentence) sandbox → dogfooded with our finance team

R (~50 words)
  - match days: 1h+ → 10–15 min (client's feedback)
  - no incorrect reconciliations in production; still in use
  - most of the time, the top suggestion was the right one
```

If the story runs over 2:10, cut the testing sentence first.

### Full text (~300 words, reference only)

> At the time, I was working on an ERP used by about five clients, and it
> had no bank reconciliation. Users had to compare their bank statement
> with the system, line by line, to make sure nothing was missing. One
> client is a soccer club, and its canteen has lots of small
> transactions. On match days, reconciling everything could take over an
> hour.
>
> My job was to fix that. I designed and built the whole feature, end to
> end, by myself, in about two months, including testing.
>
> What I did first was the import side. Transactions come from different
> sources — OFX files, PDF statements, and bank APIs through Open
> Finance — so I normalized all of them into a single internal format.
> Duplicates are blocked by a unique index in the database, scoped to
> each bank account.
>
> Then came the key trade-off: automatic matching or suggestions. For
> each transaction, I rank the five closest open bills or invoices by
> date and amount, and the user just confirms. I chose suggestions
> because a wrong automatic match is a silent error, while reviewing
> takes a couple of seconds. I only automated the predictable cases,
> like bank fees, through user-defined rules.
>
> The hardest part was that matches aren't one-to-one. One payment can
> cover several bills, and one bill can be paid in parts. So I modeled
> it as many-to-many, with a link table and a warning when the amounts
> don't add up.
>
> Before rolling it out, I tested with sandbox data and then dogfooded
> it with our own finance team.
>
> As a result, the soccer club told me match days went from over an hour
> to about ten to fifteen minutes. It's still in production, no
> incorrect reconciliations have been reported, and most of the time
> users just confirm the top suggestion.

### Follow-ups

**"Why didn't you fully automate the matching?"**
> Because the risk of an automatic match being wrong — and nobody ever
> noticing — was far more expensive than the few seconds it takes a user
> to review and confirm. I only fully automated the predictable cases,
> like bank fees, through user-defined rules.

**"How did the ranking work?"**
> Two filters first: only candidates within 10 days of the transaction,
> and only ones not reconciled yet. Then a multi-key sort: exact amount
> and date first, then the closest date, then the smallest amount
> difference. The top five are shown as suggestions.

**"What if none of the five suggestions is right?"**
> The user searches manually, by payee name or description. If the bill
> doesn't exist yet, they can create it from the transaction, and the
> system pre-fills everything it can from the bank data.

**"How did you get the bank data?"**
> Through Pluggy, an Open Finance aggregator, plus OFX and PDF uploads.
> Two banks that aren't covered by Open Finance, Asaas and Simples, were
> integrated directly through their own APIs.

**"How did you prevent duplicate transactions on re-import?"**
> Each transaction gets a deduplication key built from the source's own
> ID, prefixed by the source: the FITID for OFX, Pluggy's transaction ID,
> and for PDFs, which have no ID, a hash of the normalized fields. A
> partial unique index on (bank account, key) enforces it for active
> rows, since we use soft delete. A unique violation counts as "already
> imported", not as an error. There's also a second check by a
> fingerprint of date, type, amount and description, and possible
> duplicates are flagged for the user to confirm.
>
> *Note: re-read this part of the code myself before using these details
> in an interview.*

**"Why per bank account, and not per bank?"**
> Because the OFX spec only guarantees that a FITID is unique within an
> account, not across the whole bank.

**"How did you handle amounts that didn't add up?"**
> Matches are stored in a link table, so one transaction can point to
> several bills and the other way around. When the sums don't match, the
> system warns the user and asks if they want to correct it.

**"How did you measure the result?"**
> From the client's feedback: match days went from over an hour to about
> ten to fifteen minutes. I didn't instrument it, which is part of what
> I'd do differently.

**"What would you do differently?"**
> I'd log from day one which suggestion users actually picked, and use
> that data to replace the cascading sort with a weighted score. Right
> now, a candidate with a closer date beats one with the exact amount,
> even when the exact amount is the more likely match.

### Pronunciation

| Word | Sounds like |
|---|---|
| reconciliation | rek-un-sil-ee-AY-shun |
| reconciling | rek-un-SY-ling |
| canteen | kan-TEEN |
| statement | STAYT-ment |
| scoped | skohpt (one syllable) |
| dogfooded | DOG-foo-did |
| suggestion | sug-JES-chun |

### Practice log

| Date | Time | Missed / stuck on |
|---|---|---|
| 2026-10-20 | 2:06 | Stuck on "payables/receivables" (switched to "bills or invoices"); left out "link table" |

---

## Story #2: Teammate left out of decisions

**Answers questions like:**
- "Tell me about a time you had a conflict within your team."
- "Tell me about a time you showed leadership."
- "Tell me about a time you influenced a team without formal authority."
- "Tell me about a time you realized you were part of the problem."

**Rules specific to conflict stories:**
- The other person always has a reasonable motive. No villains.
- Show my own part in the problem.
- Resolve it directly before escalating, and keep the relationship intact.
- Spend real time on the lesson, not only on the result.
- Never name the people involved.

### Skeleton (this is what I look at)

```
S (~50 words)
  - Apple Developer Academy: 2-year program, built apps in teams
  - 2-month project, party games app, team of 5 (3 devs, 2 designers), I was a dev
  - early on: a teammate (1 of 3 devs) went quiet in decisions, ~2 weeks
    → her features, and the delivery, at risk

T (~15 words)
  - nobody asked me → I wanted to understand what was going on

A (~120 words)
  1. Talked to her one-on-one, casually
     → she felt left out: whoever spoke first decided, 2 agreeing = done
  2. MY PART: I realized I was one of them — I'd jump in over others
     → I told her that
     (short pause before this one: it's the turning point)
  3. With her okay → brought it to the team as MY proposal
     → everyone shares their view before we decide
     → explicit check that everyone agrees before closing

R (~45 words)
  - she spoke up more again, and told me it was better when I checked in
  - delivered her features; app shipped on time, on the App Store
  - (optional) prioritized must-have features → time left for design

LESSON (~35 words)
  - don't wait for someone to complain → be proactive
  - every new team after that: I made sure we agreed on HOW to decide first
```

If the story runs over 2:10, cut the optional prioritization example first.

### Full text (~270 words, reference only)

> At the time, I was at the Apple Developer Academy, a two-year program
> where we built iOS apps in small teams. On one two-month project, we
> were building a party games app. There were five of us, three
> developers and two designers, and I was one of the developers.
>
> Early on, I noticed that one of the other developers had gone quiet.
> For about two weeks, she barely spoke up when we made decisions, and
> her tasks were falling behind. With only three devs, that put her
> features, and our delivery, at risk.
>
> Nobody asked me to, but I wanted to understand what was going on. So I
> invited her for a casual one-on-one chat. She told me she felt left
> out: whoever spoke first set the direction, and as soon as two people
> agreed, the decision was made.
>
> And I realized I was part of the problem. I often jumped in with my own
> ideas over other people's. I told her that, honestly.
>
> Then, with her okay, I brought it to the team as my own proposal:
> before any decision, everyone shares their view, and we explicitly
> check that everyone agrees before closing it.
>
> As a result, she started speaking up again, and when I checked in with
> her later, she told me it felt much better. She delivered her features,
> and we shipped the app on time to the App Store.
>
> What I learned was not to wait for someone to complain before noticing
> that something's wrong. Since then, in every new team, I made sure we
> agreed on how we'd make decisions before anything else.

### Follow-ups

**"Didn't waiting for everyone to agree slow you down?"**
> A little, but we were five people in the same room, so discussions were
> short, and what we got in return was everyone committed to the
> decision. In a bigger team I wouldn't require full agreement — I'd make
> sure everyone is heard, and then let the person who owns that area
> decide, with the rest of the team committing to it even if they
> disagree.

**"What happened when you still couldn't agree?"**
> We kept discussing until we found something everyone could live with —
> usually a middle ground, or whoever proposed the idea explained their
> reasoning. For example, we couldn't fit every feature in two months, so
> we agreed on the must-haves, and that left us time to polish the
> design.

**"How did you know something was wrong?"**
> She had stopped speaking up in our decisions, and her tasks were
> falling behind. Before talking to her, I checked whether the others
> had noticed it too, so I wasn't acting on an assumption.

**"Why did you bring it to the team as your proposal, and not hers?"**
> I didn't want to put her on the spot. And the change helped everyone,
> not just her, so it made sense for it to come as a team improvement.

**"What would you do differently?"**
> I'd talk to her sooner instead of waiting two weeks. And I'd agree on
> how the team makes decisions on day one, before any problem shows up.
> That's what I started doing in every team after that.

### Pronunciation

| Word | Sounds like |
|---|---|
| quiet | KWY-et (two syllables; not "quite") |
| honestly | ON-est-lee (silent h) |
| Academy | uh-KAD-uh-mee |
| developers | dih-VEL-uh-perz |
| designers | dih-ZY-nerz (silent g) |
| explicitly | ek-SPLIS-it-lee |
| proposal | pruh-POH-zul |
| complain | kum-PLAYN |
| noticed | NOH-tist |

**"th" drill:** tongue tip between the teeth. Voiceless: *three, month,
something's*. Voiced: *the, that, there, then, other, with*. Pairs:
*three/tree, think/sink, they/day, then/den*.

### Practice log

| Date | Time | Missed / stuck on |
|---|---|---|
| 2026-10-21 | 2:05 | Nothing left out; stuck on "th" words |

---

## Story #3: The report that took the system down

**Answers questions like:**
- "Tell me about a time you failed."
- "Tell me about a mistake you made in production."
- "Tell me about a time something you built didn't scale."
- "What's a lesson you learned the hard way?"

**Rules specific to failure stories:**
- Pick a real failure with real consequences: a goal I didn't reach, not a slip.
- The sentence about the cause has "I" as the subject and my decision as the verb.
  Context (the data volume) shows up once, as a fact, never as the excuse.
- No "but" right after admitting something.
- The failure is ~30% of the answer; what I did about it and the lesson are ~70%.
- The lesson needs proof: where I applied it afterwards.

### Skeleton (this is what I look at)

```
S (~45 words)
  - only engineer on the client: audio news for radio stations
  - report: all published stories + total downloads

T (~10 words)
  - goal: a report that works AT PRODUCTION SCALE

FAILURE (~95 words)
  - "That's where I failed."
  - sorted by downloads → aggregate 11 MILLION rows (no plural on "million")
  - paginated, but page 1 still had to add up every story
  - AND I only tested with dev data (< 1,000 rows)
  - no date filter → whole system down ~5 min
  - I restarted the database (didn't know how to cancel one query)

A (~85 words): contain → fix → root cause
  1. RIGHT AWAY: told the client "it was my report" → took it offline
  2. THEN: at least one filter required (date, author or client) → back online
  3. MONTHS LATER, root cause: stored downloads column, cron every 5 min
     → sorting just reads a column
     → checked with the client that 5-minute-old numbers were fine

R (~8 words)
  - never brought the system down again   ← easy to forget: say it right after 3

LESSON (~70 words)
  - never write a list query without a limit
  - user picks at least one filter, OR there's a maximum page size
  - sort by numbers I store, not numbers I calculate every time
  - proof: law firms' ERP (hundreds of consultations/day) + Go project (default 20, cap 100)
  - next: query timeout
```

If the story runs over 2:10, cut the "checked with the client" sentence first
(it's also a follow-up below).

### Full text (~310 words, reference only)

> I was the only engineer responsible for a client that publishes audio news
> for radio stations. They asked me for a report listing every published
> story with its total number of downloads. My goal was a report that would
> work at their production scale.
>
> That's where I failed. I designed the report to sort by downloads, and that
> number came from aggregating a table with more than 11 million rows. It was
> paginated, but to build even the first page the database had to add up
> every story. And I only tested it with development data, which had fewer
> than a thousand rows. Once it was in production, someone ran it without a
> date filter, and the whole system went down for about five minutes. I
> restarted the database, because back then I didn't know I could cancel a
> single query.
>
> Right away, I told the client the problem was the report I had built, and I
> took it offline. Then I made at least one filter required, date, author or
> client, and put it back online. A few months later I went back to the root
> cause: I added a downloads column to each story, updated by a cron job
> every five minutes, so sorting just reads a column. I checked with the
> client first that five-minute-old numbers were fine for a report.
>
> It never brought the system down again. What I learned is to never write a
> list query without a limit. Either the user has to pick at least one
> filter, or there's a maximum page size. And I sort by numbers I store, not
> numbers I calculate every time. I applied the same rule in an ERP for law
> firms with hundreds of consultations a day, and my Go project has a default
> page size of 20 and a cap of 100. Next, I'm adding a query timeout there.

### Follow-ups

**"How did you find out?"**
> The client told me the system was down. I checked the database on our
> hosting panel, saw one connection stuck on that query, and restarted the
> database.

**"Why did you restart the whole database instead of killing the query?"**
> Honestly, at the time I didn't know I could cancel a single query. Today
> I'd use pg_cancel_backend or pg_terminate_backend on that process. It takes
> seconds and doesn't take the rest of the system down. Those five minutes
> of downtime were partly the cost of not knowing that.

**"Was the required filter enough?"**
> Not completely. It made the problem less likely: a filter by client alone
> could still add up millions of rows for a big client. What fixed the cause
> was storing the download count, so sorting stopped depending on the big
> table.

**"Doesn't the cron job make the numbers wrong?"**
> They're at most five minutes behind. For a downloads report that's fine.
> Nobody makes a decision based on the exact minute, and I checked that with
> the client before doing it. In exchange, the database never sums millions
> of rows to open the screen again.

**"Was there anything other people could have done differently?"**
> Not really. I was the only engineer on that client, so the design, the
> testing and the release were all mine. That's one of the reasons I picked
> this story: there's no one else to point at.

**"How would you protect against this in Go today?"**
> Three layers. Every list endpoint has a default page size and a hard cap,
> and it's always scoped to a user or a required filter. Each query gets its
> own deadline with context.WithTimeout, so the driver cancels it in the
> database if it runs too long. And the database itself has a
> statement_timeout as a safety net for any code path that forgot the
> deadline.

**"What would you do differently?"**
> Design the limit in from the start instead of adding it after the
> incident. Test list screens with a production-like volume of data, even
> if I have to generate it. And cancel the single query instead of
> restarting the database.

### Pronunciation

| Word | Sounds like |
|---|---|
| aggregating | AG-ri-gay-ting |
| paginated | PAJ-ih-nay-tid |
| consultations | kon-sul-TAY-shunz |
| column | KOL-um (silent n) |
| required | rih-KWY-erd |
| database | DAY-tuh-bays |
| cron | kron |
| million | MIL-yun (no plural before a noun: "11 million rows") |

### Practice log

| Date | Time | Missed / stuck on |
|---|---|---|
| 2026-10-21 | 1:57 | Left out the Result line ("never brought the system down again"); stuck on "every list query has to be bounded by design…" → lesson rewritten in simpler words |

---

## Story bank (not written yet)

- **B2B vs B2C split → customer research** (Apple Developer Academy). The
  team was split on the target audience for a loyalty-card app. I
  proposed interviewing potential customers, and after visiting local
  businesses the team agreed on small businesses. Fits: "Tell me about a
  time you used data to make a decision", "…convinced your team", and
  possibly #5 (disagreeing with someone).

- **Filter built in the wrong system** (ERP). A client covering elections
  asked for a filter by election round. I assumed it was for the website and
  shipped it there; it was meant for the moderation screen in the ERP, so I
  had to redo it. Fits: "Tell me about a time you misunderstood a
  requirement." Lesson: ask instead of assuming when something isn't clear.
- **Dev and production on different PHP versions** (ERP). Something worked in
  development and broke in production because I didn't check the versions
  matched. Smaller; fits "environment parity" follow-ups.
