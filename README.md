# Narwhal

Narwhal is a modern IRC bot written in Go.

![Hex.pm](https://img.shields.io/badge/irc-%23narwhal--bot%20on%20freenode-green.svg)
![Hex.pm](https://img.shields.io/badge/alpha-0.1-red.svg)
![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)

## Why

Unlike most of the world, I spend every day engaging with communities I'm a part of via IRC. This has largely been side-by-side with our robot overlords, which have up until now been fairly limited in functionality.

So, I'm building Narwhal, our tusky robot animal overlord thing, to provide functionality that did not yet exist, or was limited in scope. 

Currently this functionality includes:

- Admin / Op functionality
  - Channel-specific banning, kicking, unbanning, "unkicking". Note that Narwhal **must be an op for this**. It will try to op itself when joining a channel.
  - Add /remove users from Admins list. Supports exact user matching, host, and user@host.
  - Add / remove messages from AutoKick MessageMatches list.
  - Add / remove users from AutoKick list.
  - Add / remove users from Blacklist.
  - Ability to disable specific admin commands.
- Autokicking
  - Can automatically kick users on join or messages based on usernames and hosts, with support for basic globbing as well as regex.
  - Can automatically kick a user based on a message. Supports exact message matching, basic globbing, and regex.
  - Can automatically ban a user after they repeatedly rejoin after being kicked.
- Blacklisting
  - Ignores requests to the bot by specific users. Supports exact user matching, basic globbing, regex, and host.
- Replacer
  - Can replace own or other user messages. Uses the syntax: `.r username, s/searchword/replacerword/`
- Slap
  - Slaps the user or if itself, proclaims that it shall not.
- Song
  - Just spits out the "Narwhals Narwhals" song.
- Url Parser
  - Parse URLs to get page title.
  - Custom Reddit URL parsing to present title and score (includes upvotes)
  - Custom Youtube URL parsing to present desktop and mobile links.

## Plans

To be absolutely clear, Narwhal isn't remotely ready for use. Like I'm not even going to give you documentation on building, that's how early it is.

### General

- [x] Handle invites.
- [ ] Separate out our "plugins" into their own modules that get loaded. Move Commands -> Plugins, have an "enabled" list.

### Plugin Related

**Admin / Op functionality.**

- [x] Ability to add / remove users via IRC (must be an existing admin) for:
  - [x] Admin abilities
  - [x] Blacklist

**Autokicker:**

- [x] Expand username-based kick to support regex.
- [x] Expand on message matching.
  - [x] Should support sub-string matching and regex.

**String Replacer:**

- [x] Implement basic `.r` string replacement function. Should replace the previously stated message of the issuer or if a target, the last of the target's message.

**URL Parser:**

- [x] Implement URL parsing in messages. Beyond titles, it should specially handle:
  - [x] Reddit: Present upvotes and downvotes
  - [x] Youtube: Present mobile and non-mobile links.

## License

Narwhal is licensed under the Apache-2.0 license.