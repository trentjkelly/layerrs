# Ableton-Inspired UI Priorities

1. **Replace elevation with structural hierarchy.** Remove large shadows from track cards, menus, and dialogs. Separate areas with 1px borders, subtle panel color changes, and compact spacing instead. Start with the home feed’s `FinalTrackCard`, where `shadow-lg` currently makes each track read as a conventional web card.

2. **Adopt square surfaces everywhere.** Standardize panels, buttons, inputs, menus, artwork, and waveform bars to `border-radius: 0`. Keep a circular avatar only if the product intentionally uses it to identify artists; otherwise use square image crops.

3. **Define a restrained neutral color system plus one accent.** Use near-black for the app background, charcoal for panels, a slightly lighter neutral for hover/selected states, and off-white for primary text. Choose one saturated accent (for example orange, acid green, electric blue, or the existing violet) and reserve it for playback, selected navigation, primary actions, and focus states.

4. **Redesign tracks as dense, repeatable lanes.** Organize play control, artist, title, duration, waveform, engagement counts, and actions into consistent columns. Let the waveform be the main horizontal visual field, with a thin accent-coloured playhead/progress state. This will make the feed feel closer to music software than a social-card feed.

5. **Use a consistent 4px or 8px layout grid.** Reduce oversized header/card padding and align the sidebar, header, track metadata, actions, and waveform to shared vertical columns. Consistent density will make the app feel more precise and professional.

6. **Simplify typography into a functional hierarchy.** Use `Host Grotesk` for the interface, retaining Humankind only for the Layerrs wordmark if desired. Avoid `DM Serif Display` in product UI. Prefer concise medium-weight labels and small uppercase metadata labels such as `USES`, `APPEARS ON`, and `PLAYS`.

7. **Make controls feel crisp rather than floaty.** Replace scale and shadow hover effects with quick background changes, outline changes, or accent-coloured text. Keep transitions brief and avoid animation on controls that should feel like instrument controls.

8. **Apply the same visual rules to all secondary screens.** Update page forms, upload, profile, library, and modal views after the home feed is established. This includes removing remaining rounded inputs/buttons and replacing mixed default blue, green, red, indigo, and violet action colors with the chosen semantic/accent palette.
