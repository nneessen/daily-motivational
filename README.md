# Morning texts

Every weekday at 5:17 AM Denver time, a cron-job.org timer tells GitHub to run this program. It emails you (the `TO_ADDRESS` secret) 4 options for the morning text to the team. Each one is a real quote from a real person with one short line from you under it. Pick one, tweak it, send it.

The quotes come from `quotes.txt` and are printed exactly as written there. Claude only writes the short line under each one, so it can't misquote anybody or make a quote up. Each weekday takes the next 4 quotes, so none comes back until the whole list has been used. Adding or removing a quote is fine: at most a quote or two from recent mornings comes back early.

## Change what it writes
- **Quotes:** edit `quotes.txt`, one per line: `the quote | who said it`. Check who really said it at quoteinvestigator.com or en.wikiquote.org before adding one. A lot of famous quotes are pinned on the wrong person.
- **The line under each quote:** edit `voice.md`. It's the prompt, with your real example message in it.
- **This week's news:** write a line or two in `focus.txt` (like "Contest ends Friday."). It may show up in one of the lines. Clear it out when it's old. The repo is public, so only write what you'd be fine with anyone reading.

Edits on github.com count too. The next morning's run picks them up.

## One-time setup
1. **Keys:** both are the ones commissionTracker already uses.
   - The Anthropic key costs about $1–3 a month here.
   - Mailgun sends one email a day from `morning@updates.thestandardhq.com`.
2. **Local testing:** `cp .env.example .env` and fill in the two keys and `TO_ADDRESS`.
3. **GitHub secrets:** `gh secret set -f .env` loads all four from `.env`: `ANTHROPIC_API_KEY`, `MAILGUN_API_KEY`, `MAILGUN_DOMAIN` and `TO_ADDRESS`.

## Run it
```
scripts/test_run.sh          # checks + prints today's email, sends nothing
go run . -dry-run -date 2026-11-02   # print the email for another day
scripts/test_run.sh --send   # checks + sends one real email
gh workflow run daily.yml    # run the GitHub job now
```

## Schedule
The timer lives at cron-job.org, set to America/Denver, so daylight saving time takes care of itself. At 5:17 AM Monday to Friday it sends `POST https://api.github.com/repos/nneessen/daily-motivational/actions/workflows/daily.yml/dispatches` with body `{"ref":"main"}`. It uses a fine-grained GitHub token that can only run this one repo's Actions. cron-job.org emails you if a call fails. When the token expires, make a new one and paste it into the job's `Authorization` header.

There's no GitHub schedule (`on: schedule`) in the workflow on purpose. GitHub never fired it for this repo during its Oct 21–22, 2026 Actions outages. If one were added back while the timer runs too, you'd get two emails every morning.

The repo is public because public repos get unlimited free Actions minutes, and theStandardHQ uses up this account's private-repo minutes every month. The keys live in GitHub secrets, which stay hidden in a public repo.
