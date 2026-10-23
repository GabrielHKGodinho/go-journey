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
- **Pause instead of "So".** Starting every sentence with "So" sounds
  like hesitation; a half-second pause sounds confident.
- **Details are follow-ups, not part of the story.** Keep the story at the
  level of decisions; each detail left out is a follow-up where I can
  go deep.

| # | Story | Main competencies | Status |
|---|---|---|---|
| 1 | Bank reconciliation (ERP) | Technical challenge, ownership, trade-offs | Done |
| 2 | Teammate left out of decisions (Apple Developer Academy) | Conflict, leadership without authority, self-awareness | Done |
| 3 | The report that took the system down (ERP) | Accountability, learning from failure, scaling | Done |
| 4 | Year-end reports on a tight deadline (ERP) | Prioritization, scope negotiation, customer focus | Done |
| 5 | Polling vs. webhook (ERP) | Disagreeing with a manager, influence, designing for failure | Done |
| 6 | Production passwords out of git (ERP) | Ownership, security, process improvement | Done |

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

## Story #4: Year-end reports on a tight deadline

**Answers questions like:**
- "Tell me about a time you worked under a tight deadline."
- "Tell me about a time you had to prioritize."
- "Tell me about a time you pushed back on a requirement."
- "Tell me about a time you found a simpler solution than what was asked."

**Rules specific to deadline stories:**
- The action is a decision about scope, not extra hours.
- The Situation has to show *why* the deadline was tight.
- Say who I told, and when.

### Skeleton (this is what I look at)

```
S (~45 words)
  - law firm client (ERP), early December 2024
  - asked for a reports feature: 5 years of service records + courses
  - for their year-end conference in January, and they wanted to review it first
    → ~15 business days, with the holidays
  - as a feature: lots of data, many screens and filters
    + I had open tickets from other clients

T (~20 words)
  - my job: get them what they needed in time,
    without leaving my other clients waiting

A (~130 words)
  1. KEY INSIGHT: the conference needed the NUMBERS, not a screen
     → I proposed: manual report now, the feature next month
     → explained why on a call; client agreed — they just needed the data
  2. SQL queries on production (off-hours) → document in THEIR layout → PDF
     ~1 week
  3. Accuracy: checked small samples myself + client reviewed before the conference

R (~45 words)
  - delivered on time, reviewed, used at the conference
  - next month: shipped the feature, as agreed
  - BUILT ON THE SAME QUERIES the client had already validated
    → the workaround made the real feature safer, not wasted work
```

If the story runs over 2:10, cut "(off-hours)" first (it's also a follow-up below).

### Full text (~300 words, reference only)

> At the time, I was working with a law firm client on their ERP. In early
> December 2024, they asked for a reports feature: five years of service
> records plus courses, for their year-end conference in January, and they
> wanted to review it first. That gave me about 15 business days because of
> the holidays, and the feature involved lots of data, many screens and
> filters. It also wasn't my only task — I had open tickets from other
> clients.
>
> My goal was to get them what they needed in time, without leaving my other
> clients waiting.
>
> The key thing I realized was that for the conference they needed the
> numbers, not a new screen. So I proposed a manual report now and the
> feature next month, and explained why on a call. The client agreed — they
> just needed the data. I wrote SQL queries on production, off-hours because
> of the volume, put the results into a document in their layout and
> exported it as a PDF. That took about a week. To make sure the numbers were
> right, I checked small samples myself and sent it to the client to review
> before the conference.
>
> In the end, it was delivered on time, reviewed, and they used it at the
> conference. The next month I shipped the feature as agreed — and the best
> part was that it was built on the same queries the client had already
> validated, so the workaround wasn't wasted work; it made the feature safer.

### Follow-ups

**"Why not just build the feature faster, or work extra hours?"**
> Because the feature wasn't what they needed for the conference — the
> numbers were. Rushing many screens and filters in two weeks, with other
> clients waiting, would risk bugs in exactly the data they were going to
> present. The manual report delivered what mattered, with less risk.

**"You ran queries directly on production?"**
> Yes, off-hours, because they went over five years of data. I'd already
> learned that lesson the hard way — a report of mine once took a system
> down — so I made sure they ran when nobody was using it.

**"How did you know the numbers were right?"**
> Two checks. I compared small samples against records I could verify in
> the system, and I sent the report to the client to review before the
> conference. They knew their numbers better than anyone.

**"How did the client feel about not getting the feature on time?"**
> They took it well. I explained why on a call, and they only needed the
> data for the conference. The feature came the next month, as we'd agreed.

**"What would you do differently?"**
> Most of that week went into building the document in their layout. I'd
> send the raw numbers in a spreadsheet first, for the client to validate,
> and only then do the layout. If a number was wrong, we'd find out before
> the formatting was done.

**To prepare:** how I handled the other clients' tickets that week (the Task
says "without leaving my other clients waiting" — be ready to back it up).

### Pronunciation

| Word | Sounds like |
|---|---|
| queries | KWEER-eez (not "kerries") |
| business days | BIZ-ness days (two syllables) |
| workaround | WURK-a-round |
| conference | KON-fer-ence (stress on the first syllable) |
| layout | LAY-out |
| accuracy | AK-yur-uh-see |
| validated | VAL-ih-day-tid |

### Practice log

| Date | Time | Missed / stuck on |
|---|---|---|
| 2026-10-23 | 1:50 | Nothing left out; "they thought it was quicker to develop the feature" said the opposite (cut it); "secure" → "safer"; working **with** a client; agree (no object); 8 sentences starting with "So" |

---

## Story #5: Polling vs. webhook

**Answers questions like:**
- "Tell me about a time you disagreed with your manager."
- "Tell me about a time you convinced someone to change their mind."
- "Tell me about a technical decision you pushed back on."
- "Tell me about a time you designed for failure."

**Rules specific to disagreement stories:**
- Show the other person's side first. Their reason was fair.
- Convince with something concrete (an example, a cost, a risk), never by insisting.
- Answer their concern, don't just win the argument.

### Skeleton (this is what I look at)

```
S (~40 words)
  - this year, same law firm client (ERP)
  - consultants open a service record for every phone call (200+ a day) — by hand
    → calls could be missed; goal: automate it
  - my manager's plan: POLL the phone provider's API

T (~20 words)
  - I was building it, and I believed polling was the wrong design
    → my job: convince him, or be convinced

A (~140 words)
  1. FIRST, HIS SIDE: polling is simpler and doesn't depend on the provider
     — a fair point
  2. MY CONCERNS:
     - load: most requests come back with "nothing changed"
     - delay: a short call could end before the next check,
       and consultants would wait for it to start the record
  3. PROPOSAL: webhook for real time + slow polling (every 10 min) as a SAFETY NET
     → showed an example of the webhook code
  4. He pushed back: "it depends on them" → I walked through load + delay
     + the safety net COVERED HIS CONCERN → he agreed
  5. I emailed the provider → they agreed, ~2 weeks on their side;
     I built the endpoint (token auth)

R (~50 words)
  - the call shows up on the consultant's screen as they answer it
  - in production since then, still in use; client liked it
  - missed by the webhook → pending-calls screen + warning → never lost silently
  - LESSON: disagree with something concrete, and address the other person's concern
```

If the story runs over 2:10, cut "(token auth)" first, then "~2 weeks on their side".
If #4 was told in the same interview, open with "With that same law firm client…".

### Full text (~330 words, reference only)

> Earlier this year, with that same law firm client, the consultants had to
> open a service record by hand for every phone call they received — over
> 200 calls a day. Doing it by hand meant calls could be missed, so the goal
> was to automate it. My manager's plan was to poll the phone provider's API.
>
> I was the one building it, and I believed polling was the wrong design, so
> my job was to convince him — or be convinced.
>
> First, I tried to understand his side: polling is simpler and doesn't
> depend on the provider, and that's a fair point. But I had two concerns.
> Most requests would come back with nothing changed, wasting server
> resources. And there's a delay: a short call could end before the next
> check, and consultants would have to wait for it before they could start
> filling in the record.
>
> What I proposed was a webhook for real time, plus a slow polling job every
> ten minutes as a safety net, so any call the webhook missed would still be
> caught. I showed him an example of the webhook code. He pushed back: "it
> depends on them." I walked him through the load and the delay again, and
> showed him that the safety net already covered his concern — and he
> agreed. I emailed the provider, they agreed too, and it took them about two
> weeks on their side. Meanwhile, I built the endpoint, with token
> authentication.
>
> In the end, the call shows up on the consultant's screen as they answer the
> phone. It's been in production since then, and the client said they liked
> it. Any call the webhook misses goes to a pending-calls screen, with a
> warning on the main screen, so no call is ever lost silently. What I took
> from it is to disagree with something concrete, and to address the other
> person's concern.

### Follow-ups

**"So you ended up polling anyway?"**
> Yes, but a very different kind. To be close to real time, polling would
> have to run every few seconds. The safety net runs every ten minutes —
> six requests an hour — because it doesn't need to be fast. Real time comes
> from the webhook; the polling only makes sure nothing gets lost.

**"How do you make sure a webhook call really comes from the provider?"**
> The provider sends a token with every call, and the endpoint rejects any
> request without the right one. If the provider supported it, I'd also
> verify a signature of the request body (HMAC), so a leaked token alone
> wouldn't be enough.

**"What if your manager had still said no?"**
> Then I'd build the polling, and build it well. It was his call to make,
> and his concern was legitimate. I'd probably suggest we measure the load
> and the delay after a few weeks, so we could revisit it with real numbers.

**"How did it affect your relationship with your manager?"**
> It stayed the same. It was a technical discussion, not a personal one, and
> we kept working together normally.

**"What would you do differently?"**
> From what consultants see, most calls arrive through the webhook, but we
> don't measure it. I'd track how many calls only arrive through the safety
> net, so we'd have a number for how reliable the webhook is, instead of a
> feeling.

**To prepare:** what happens if the provider sends the same webhook twice
(does the system create two records, or is the call ID unique?).

### Pronunciation

| Word | Sounds like |
|---|---|
| polling | POH-ling ("poll", not "pull") |
| webhook | WEB-hook (the h is pronounced) |
| provider | pro-VY-der |
| concerns | kon-SURNZ |
| pushed back | PUSHT back |
| reliable | ri-LY-uh-bl |
| endpoint | END-point |
| safety net | SAYF-tee net (two syllables in "safety") |
| caught | kawt (past of "catch") |

### Practice log

| Date | Time | Missed / stuck on |
|---|---|---|
| 2026-10-23 | 2:01 | Nothing left out; added the consultant-waiting argument (kept); "his side was" → "I tried to understand his side"; tangled the "safety net covered his concern" line; "is in production since" → "has been in production since then"; 12 "So" |

---

## Story #6: Production passwords out of git

**Answers questions like:**
- "Tell me about a time you took ownership of a problem."
- "Tell me about something you fixed that wasn't your job."
- "Tell me about a time you improved a process."
- "Tell me about a security issue you found."

**Rules specific to ownership stories:**
- Open the Task with "nobody asked me".
- Show that I asked before changing something other people depended on.
- Don't claim the other projects adopted it because of me — I don't know that.
- Say the gap (no rotation) before the interviewer finds it.

### Skeleton (this is what I look at)

```
S (~50 words)
  - 2022, I was an intern, ~1 year in; ERP project for one client
  - production passwords (database, APIs) were committed to the GitHub repo
    → anyone with access to the repo had the keys to production
      (and so would anyone who got into one of those laptops or accounts)
  - deploys: manual FTP upload, overwriting config files
    → it had ALREADY HAPPENED, more than once: dev config uploaded over production

T (~20 words)
  - nobody had asked me to fix it — but it was a security risk
    and a mistake that kept happening → I took it on

A (~110 words)
  1. FIRST, ASKED MY MANAGER: explained the risk + a better way → approved
  2. config out of the code:
     - one .env per environment (phpdotenv)
     - .env in .gitignore + .env.example (keys, no values)
     → an FTP upload can't overwrite the environment's config anymore
  3. removed the passwords from the git history
     → so they wouldn't keep showing up FROM THEN ON

R (~50 words)
  - no more deploys with the wrong environment's credentials
    → deploying got easier and safer, one place per environment
  - later, the other projects adopted the same setup;
    new projects start with it
  - LOOKING BACK: I should have rotated the passwords
    + automated deploys (GitHub Actions)
```

If the story runs over 2:10, cut "+ automated deploys" first.

### Full text (~290 words, reference only)

> Back in 2022, I was an intern, about a year into the company, working on an
> ERP project for one client. I noticed that their production passwords —
> database and APIs — were committed to the GitHub repository. Anyone with
> access to the repo had the keys to production, and so would anyone who got
> into one of those laptops or GitHub accounts. On top of that, deploys were
> manual FTP uploads, and it had happened more than once that someone
> uploaded a file with the development passwords over production.
>
> Nobody had asked me to fix it, but it was a security risk and a mistake
> that kept happening, so I took it on.
>
> First, I asked my manager. I explained the risk and a better way to do it,
> and he approved it right away. Then I took the config out of the code: one
> .env file per environment, using the phpdotenv library, with .env added to
> the .gitignore, and a .env.example so anyone could see which keys they
> needed. I also removed the passwords from the git history, so they wouldn't
> keep showing up for anyone cloning the repo from then on.
>
> In the end, we had no more deploys with the wrong environment's
> credentials, and deploying got easier and safer, with all the passwords in
> one place per environment. Later, the other projects adopted the same
> setup, and new projects start with it. Looking back, I should have rotated
> the passwords. Anyone who had cloned the repo, or simply copied them
> somewhere, could still use them, so it would have been safer. I'd also
> have automated deploys with GitHub Actions.

### Follow-ups

**"You removed the passwords from git history. Isn't that enough?"**
> No. Cleaning the history only stops people from finding the passwords
> from now on. Anyone who had already cloned the repo still has them. The
> real fix is to rotate the passwords, and I didn't do that at the time.
> Today, I'd do it first.

**"Did the other projects adopt it because of you?"**
> I believe my manager took it from my project, but I can't say for sure.
> What I know is that the other projects use the same setup now, and new
> projects start with it.

**"Why did you ask your manager instead of just doing it?"**
> It changed how the production config and the deploys worked, and other
> people deployed that project too. They needed to know, and as an intern it
> wasn't a call I should make alone.

**"How do the passwords get to production now?"**
> The production .env lives on the server, and it's set up by hand once.
> It's not in the repo, so a deploy upload never touches it.

**"How would you do it in Go today?"**
> The same idea: config comes from environment variables, read at startup,
> and the service refuses to start if a required one is missing. In my Go
> project, the API returns an error right away if DATABASE_URL or
> RABBITMQ_URL is empty. Secrets live in the platform (the hosting
> provider's variables, GitHub Actions secrets), never in the repo.

**"What would you do differently?"**
> Rotate the passwords right after removing them, and automate the deploy
> with GitHub Actions, so nobody uploads config files by hand at all.

**To prepare:** how exactly I removed the passwords from the history (git
filter-repo, BFG, or something else) — check before the next mock.

### Pronunciation

| Word | Sounds like |
|---|---|
| .env | dot-EN-vee |
| .gitignore | dot-GIT-ig-nor |
| phpdotenv | P-H-P dot env |
| credentials | kreh-DEN-shulz |
| intern | IN-turn |
| rotated | ROH-tay-tid |
| environment | en-VY-ron-ment |
| repository | rih-POZ-ih-tor-ee |

### Practice log

| Date | Time | Missed / stuck on |
|---|---|---|
| 2026-10-23 | 1:58 | Nothing left out; stuck on the risk sentence (sounded like accusing the client's IT managers → rewritten); "nobody could see past passwords" contradicted the lesson → "from then on"; could **steal**; "it was safer" → "it would have been safer"; 9 "So" |

---

## Story bank (not written yet)

- **B2B vs B2C split → customer research** (Apple Developer Academy). The
  team was split on the target audience for a loyalty-card app. I
  proposed interviewing potential customers, and after visiting local
  businesses the team agreed on small businesses. Fits: "Tell me about a
  time you used data to make a decision" and "…convinced your team".

- **Filter built in the wrong system** (ERP). A client covering elections
  asked for a filter by election round. I assumed it was for the website and
  shipped it there; it was meant for the moderation screen in the ERP, so I
  had to redo it. Fits: "Tell me about a time you misunderstood a
  requirement." Lesson: ask instead of assuming when something isn't clear.
- **Dev and production on different PHP versions** (ERP). Something worked in
  development and broke in production because I didn't check the versions
  matched. Smaller; fits "environment parity" follow-ups.
- **Word integration for the year end** (law firm ERP, a different December
  from #4). The client needed Word documents generated from the system before
  the new year; I prioritized what made it work and negotiated the
  improvements for later. Backup for #4.
- **Hardcoded options → managed by the client** (radio client). Options like
  story areas were hardcoded, so every new one was a request to us. Nobody
  asked me; I moved them to the database with a screen for the client to
  manage them. Fits: "reduced repetitive work", "thought about the customer";
  backup for #6.

**Gap to fill (not a story yet):** "Tell me about a time you had to learn
something new quickly." This one will come up because of the PHP → Go move.
