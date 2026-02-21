#!/usr/bin/env python3
"""Mark proper nouns in world_measurements.json and world_durations.json.

Proper nouns are unique, named things where you'd say "the X" not "a X".
Examples: "the Eiffel Tower", "the Grand Canyon", "the Moon"
Counter-examples: "a Soccer Ball", "an African Elephant", "a School Bus"
"""

import json
from pathlib import Path

DATA_DIR = Path(__file__).resolve().parent.parent.parent / "internal" / "data"

# ── Proper nouns in world_measurements.json ─────────────────────────────────

# Categories where ALL items are proper nouns
ALWAYS_PROPER_CATEGORIES = {"Country", "City", "Celestial", "Artifact"}

# Explicit proper nouns in other categories
MEASUREMENTS_PROPER_NOUNS = {
    # Structure (named landmarks, unique buildings)
    "Eiffel Tower",
    "Burj Khalifa",
    "Empire State Building",
    "Sydney Opera House",
    "Taj Mahal",
    "Colosseum Rome",
    "Great Wall of China (total)",
    "Golden Gate Bridge",
    "Brooklyn Bridge",
    "CN Tower",
    "Statue of Liberty",
    "Christ the Redeemer",
    "Great Pyramid of Giza",
    "Stonehenge",
    "Pantheon Rome",
    "Notre-Dame Cathedral",
    "Big Ben (Elizabeth Tower)",
    "Leaning Tower of Pisa",
    "Washington Monument",
    "Hoover Dam",
    "Three Gorges Dam",
    "Panama Canal",
    "Channel Tunnel",
    "Burj Al Arab Hotel",
    "One World Trade Center",
    "Shanghai Tower",
    "Petronas Towers",
    "Space Needle Seattle",
    "Arc de Triomphe",
    "Brandenburg Gate",
    "Nelson's Column",
    "London Eye",
    "Dubai Ferris Wheel (Ain Dubai)",
    "Trans-Siberian Railway",
    "Suez Canal",
    "Akashi Kaikyo Bridge",
    "George Washington Bridge",
    "London Tower Bridge",
    "Sydney Harbour Bridge",
    "Millau Viaduct",
    "Rialto Bridge Venice",
    "Charles Bridge Prague",
    "Trevi Fountain",
    "Flatiron Building",
    "Chrysler Building",
    "Willis Tower (Sears Tower)",
    "Transamerica Pyramid",
    "Taipei 101",
    "Shard (London)",
    "Gherkin (London)",
    "Sagrada Familia",
    "Hagia Sophia",
    "St. Peter's Basilica",
    "Westminster Abbey",
    "Angkor Wat",
    "Machu Picchu (complex)",
    "Chichen Itza (El Castillo)",
    "Parthenon Athens",
    "Acropolis (platform)",
    "Alhambra Palace (complex)",
    "Versailles Palace",
    "Buckingham Palace",
    "White House",
    "US Capitol Building",
    "Pentagon",
    "Vatican (total area)",
    "Kremlin (total complex)",
    "Forbidden City",
    "Olympic Stadium (Bird's Nest)",
    "Sydney Olympic Stadium",
    "Louvre Museum",
    "Metropolitan Museum of Art",
    "British Museum",
    "Vatican Museums",
    "Hermitage Museum",
    "National Mall (Washington DC)",
    "Times Square",
    "Red Square Moscow",
    "Tiananmen Square",
    "Trafalgar Square",
    "St. Peter's Square Rome",

    # Natural Feature (unique named features)
    "Mount Everest",
    "Grand Canyon",
    "Amazon River",
    "Nile River",
    "Mariana Trench",
    "Victoria Falls",
    "Angel Falls",
    "Niagara Falls",
    "Great Barrier Reef",
    "Sahara Desert",
    "Antarctica",
    "Lake Superior",
    "Caspian Sea",
    "Atlantic Ocean",
    "Pacific Ocean",
    "Dead Sea",
    "Greenland Ice Sheet",
    "Congo River",
    "Mississippi River",
    "Andes Mountains (length)",
    "Himalayan Range",
    "Asteroid (Ceres)",
    "Mount Kilimanjaro",
    "Mount Fuji",
    "Yellowstone Caldera",
    "Great Blue Hole Belize",
    "Amazon Rainforest",
    "Siberian Taiga Forest",
    "Gibraltar Rock",
    "Krakatoa Volcano",
    "Vesuvius Volcano",
    "Mount St. Helens",
    "Mississippi Delta",
    "Okefenokee Swamp",
    "Lake Baikal",
    "Nile Delta",
    "Ayers Rock (Uluru)",
    "Table Mountain",
    "Mount Cook (New Zealand)",
    "K2",
    "Aconcagua",
    "Mont Blanc",
    "Matterhorn",
    "Ben Nevis",
    "Snowdon",
    "Inca Trail",
    "Appalachian Trail",
    "Pacific Crest Trail",
    "Rhine River",
    "Danube River",
    "Yangtze River",
    "Ganges River",
    "Colorado River",
    "Hudson River",
    "Thames River",
    "Seine River",
    "Volga River",
    "Lake Victoria",
    "Lake Huron",
    "Lake Michigan",
    "Lake Erie",
    "Lake Ontario",
    "Lake Titicaca",
    "Red Sea",
    "Mediterranean Sea",
    "North Sea",
    "Caribbean Sea",
    "Arctic Ocean",
    "Central Park NYC",
    "Hyde Park London",

    # Spacecraft (unique named spacecraft)
    "Saturn V Rocket",
    "SpaceX Falcon 9",
    "SpaceX Starship",
    "Space Shuttle (stack)",
    "International Space Station",
    "Voyager 1 Probe",
    "Mars Curiosity Rover",
    "Hubble Space Telescope",
    "James Webb Space Telescope",
    "Ariane 5 Rocket",
    "Mars Perseverance Rover",

    # Aircraft (unique named aircraft)
    "Antonov An-225",
    "Dirigible (Hindenburg)",

    # Watercraft (specific named ships)
    "Titanic (RMS)",
    "Aircraft Carrier (USS Gerald R. Ford)",
    "Cruise Ship (Symphony of the Seas)",
    "Container Ship (Emma Maersk)",
    "Battleship (USS Missouri)",
    "Destroyer (Arleigh Burke class)",
    "Aircraft Carrier (Nimitz class)",
    "Submarine (Ohio class)",

    # Equipment (unique facilities)
    "Large Hadron Collider",
    "Radio Telescope (Arecibo, diameter)",
    "Telescope (Yerkes 40-inch refractor)",

    # Object (unique named objects)
    "Big Ben Bell (Great Bell)",
}

# ── Proper nouns in world_durations.json ────────────────────────────────────

DURATIONS_PROPER_NOUNS = {
    # Historical events (specific named events)
    "Apollo 11 Moon landing (descent)",
    "Apollo 11 moonwalk (first EVA)",
    "Apollo 11 total mission",
    "First powered airplane flight (Wright 1903)",
    "Titanic sinking",
    "Hiroshima bomb detonation to shockwave",
    "World War I duration",
    "World War II duration",
    "Cold War duration",
    "Roman Empire duration (Western)",
    "British Empire at peak (duration)",
    "Hundred Years' War",
    "Thirty Years' War",
    "American Civil War",
    "French Revolution",
    "Siege of Leningrad",
    "Chernobyl explosion to reactor fire extinguished",
    "Berlin Wall standing",
    "D-Day Normandy invasion (June 6, 1944)",
    "Battle of Gettysburg",
    "Cuban Missile Crisis",
    "Moon race (Sputnik to Apollo 11)",
    "Construction of Eiffel Tower",
    "Construction of Great Pyramid of Giza",
    "Black Death pandemic (Europe)",
    "Spanish Flu pandemic 1918",
    "COVID-19 pandemic (declared to end of PHE)",
    "Great Fire of London 1666",
    "Pompeii eruption (Vesuvius 79 AD)",
    "Construction of Panama Canal",
    "Construction of Empire State Building",
    "Manhattan Project (Trinity to Hiroshima)",
    "Space Shuttle Challenger disaster",
    "Space Shuttle Columbia reentry disaster",
    "Voyager 1 reaching interstellar space",
    "Wright Brothers first flight to Moon landing",
    "First iPhone to present (2007–2025)",
    "Agricultural revolution",
    "Industrial Revolution duration",
    "First solo nonstop transatlantic flight (Lindbergh)",
    "Apollo 13 crisis duration",
    "Shortest war in history (Anglo-Zanzibar)",

    # Geology (specific named events)
    "2004 Indian Ocean Tsunami (wave travel)",
    "2011 Tōhoku earthquake duration",
    "1906 San Francisco earthquake",
    "1980 Mount St. Helens eruption (initial blast)",
    "1883 Krakatoa eruption",
    "Chicxulub asteroid impact",
    "Ice Age (last glacial maximum)",
    "Formation of Grand Canyon",
    "Tambora eruption 1815",
    "Formation of Hawaiian islands",
    "K-Pg extinction event",

    # Weather (specific named events)
    "1925 Tri-State Tornado",
    "Hurricane Katrina (Cat 5 peak duration)",
    "Great Blizzard of 1888",
    "1815 'Year Without a Summer'",
    "Dust Bowl period",
    "Little Ice Age",
    "London Great Smog 1952",

    # Astronomy (unique phenomena)
    "Big Bang to first stars",
    "Age of the Universe",
    "Formation of Solar System",
    "Age of Earth",
    "Sun's remaining lifespan (main sequence)",
    "Jupiter's rotation (one day)",
    "Venus rotation (one day)",
    "Mars rotation (one day)",
    "Saturn's orbit (one year)",
    "Pluto's orbit (one year)",
    "Light travel from nearest star (Proxima Centauri)",
    "Light travel from Andromeda Galaxy",
    "Halley's Comet orbital period",
    "Voyager 1 to reach nearest star (theoretical)",
    "Milky Way galactic rotation",
    "ISS orbital period",

    # Sports (specific events)
    "Tour de France (average)",
    "Longest tennis match (Isner–Mahut 2010)",
    "Longest spacewalk (March 2001)",

    # Culture
    "Beethoven's 9th Symphony",
    "Oscar ceremony (average)",
    "Super Bowl game (average)",

    # Biology (specific named events)
    "Dinosaur extinction to humans",
    "Evolution of Homo sapiens",
    "First life on Earth to present",
    "World's longest recorded hiccup bout",
}


def mark_proper_nouns(path, proper_set, use_category_rule=False):
    with open(path) as f:
        items = json.load(f)

    marked = 0
    for item in items:
        is_proper = False
        if use_category_rule and item.get("category") in ALWAYS_PROPER_CATEGORIES:
            is_proper = True
        if item["name"] in proper_set:
            is_proper = True

        if is_proper:
            item["proper_noun"] = True
            marked += 1
        else:
            item.pop("proper_noun", None)

    with open(path, "w") as f:
        json.dump(items, f, indent=2, ensure_ascii=False)
        f.write("\n")

    return len(items), marked


def verify(path, proper_set, use_category_rule=False):
    """Check for items in proper_set that don't exist in the file."""
    with open(path) as f:
        items = json.load(f)
    names = {item["name"] for item in items}
    missing = proper_set - names
    if missing:
        print(f"  WARNING: {len(missing)} names in proper set not found in data:")
        for name in sorted(missing):
            print(f"    - {name}")


def main():
    m_path = DATA_DIR / "world_measurements.json"
    d_path = DATA_DIR / "world_durations.json"

    print("Marking proper nouns in world_measurements.json...")
    verify(m_path, MEASUREMENTS_PROPER_NOUNS, use_category_rule=True)
    total, marked = mark_proper_nouns(m_path, MEASUREMENTS_PROPER_NOUNS, use_category_rule=True)
    print(f"  {marked}/{total} items marked as proper nouns")

    print("\nMarking proper nouns in world_durations.json...")
    verify(d_path, DURATIONS_PROPER_NOUNS, use_category_rule=False)
    total, marked = mark_proper_nouns(d_path, DURATIONS_PROPER_NOUNS, use_category_rule=False)
    print(f"  {marked}/{total} items marked as proper nouns")


if __name__ == "__main__":
    main()
