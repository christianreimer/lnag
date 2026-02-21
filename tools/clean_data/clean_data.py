#!/usr/bin/env python3
"""Clean world_measurements.json and world_durations.json.

Removes items where:
1. The measurement is not commonly known by people
2. The name is derived from an SI or imperial unit
3. There are duplicates (keeps the more generic version)

Also renames unit-derived names to generic versions.
"""

import json
import sys
from pathlib import Path

DATA_DIR = Path(__file__).resolve().parent.parent.parent / "internal" / "data"

# ── Items to REMOVE from world_measurements.json ────────────────────────────

MEASUREMENTS_REMOVE = {
    # ── Unit-derived names (the item IS a unit or is defined by one) ──
    "Gallon Jug of Water",
    "50-gallon Drum",
    "Garden Hose (50ft)",
    "Extension Cord (25ft)",
    "Propane Tank (20lb)",
    "Bag of Potatoes (50lb)",
    "Sack of Flour (25kg)",
    "Paint Can (1 gallon)",
    "Dumbbell (20lb)",
    "Watering Can (2 gallon)",
    "Ruler (30cm)",
    "Yardstick",
    "Tank (propane, 1000 gallon)",
    "Extension Ladder (32ft)",  # duplicate of Ladder (extension, 24ft) after rename

    # ── Duplicates (keep the more generic/better version) ──
    "ATM Machine",              # keep "ATM"
    "ATM (bank)",               # keep "ATM"
    "Liberty Bell (replica)",   # keep "Liberty Bell"
    "Concert Grand Piano (Steinway D)",  # keep "Grand Piano"
    "Baby Grand Piano",         # keep "Grand Piano"
    "Giant Manta Ray (full span)",  # keep "Giant Manta Ray"
    "Colosseum (height)",       # keep "Colosseum Rome"
    "Dead Sea (depth)",         # keep "Dead Sea"
    "Mariana Trench (deepest point)",  # keep "Mariana Trench"
    "CERN Hadron Collider Tunnel",  # keep "Large Hadron Collider (circumference)"
    "Freedom Tower (NYC)",      # duplicate of "One World Trade Center"
    "Caspian Sea depth",        # keep "Caspian Sea"
    "Swimming Pool (Olympic, volume)",  # keep "Olympic Swimming Pool"
    "Keg of Beer (half barrel)",  # keep "Beer Keg (standard)"
    "Golf Buggy",               # duplicate of "Golf Cart"
    "Boxing Ring (standard)",   # keep "Boxing Ring"
    "Surfboard (longboard)",    # keep "Surfboard" (in Watercraft)
    "Brandenburg Gate Column",  # "Brandenburg Gate" has more data
    "Escalator",                # keep "Escalator (standard)" which is more realistic size
    "Billiard Table (standard)",  # keep "Pool Table" after rename
    "Ship's Anchor (large)",    # keep "Anchor (ship)"

    # ── Not commonly known measurements ──

    # Personal/grooming (nobody knows these measurements)
    "Toothbrush",
    "Toothpaste Tube",
    "Bar of Soap",
    "Shampoo Bottle (standard)",
    "Razor (safety)",
    "Hairbrush",
    "Comb",
    "Nail Clippers",

    # Tiny fasteners/stationery (measurements meaningless to people)
    "Paper Clip",
    "Safety Pin",
    "Button (shirt)",
    "Zipper (jacket length)",
    "Rubber Band",
    "Luggage Tag",
    "Letter Envelope",

    # Small accessories (obscure measurements)
    "Watch Band",
    "Ring (wedding band)",
    "Earring",
    "Bracelet",
    "Necklace (standard)",

    # Medical devices/items (obscure measurements)
    "Hearing Aid",
    "Contact Lens",
    "Dentures (full set)",
    "Hip Replacement (implant)",
    "Pacemaker",
    "Respirator Mask (N95)",
    "Surgical Gloves (pair)",
    "Band-Aid (standard)",
    "Syringe (10mL)",
    "Aspirin Tablet",
    "Prescription Pill Bottle",
    "Cast (leg plaster)",
    "IV Drip Stand",
    "Oxygen Tank (standard)",
    "Walker (medical)",
    "Crutch",
    "Gurney (stretcher)",

    # Obscure tools (specialized, nobody knows measurements)
    "Wire Stripper",
    "Voltmeter",
    "Multimeter",
    "Soldering Iron",
    "Heat Gun",
    "Pipe Cutter",
    "Cable Cutter (large)",
    "Pipe Wrench (large)",
    "Torque Wrench",
    "Belt Sander",
    "Angle Grinder",
    "Drill Press",
    "Bench Vise",
    "Square (carpenter's)",
    "Compass (drafting)",
    "Protractor",
    "Level (carpenter's 4ft)",
    "Compass (navigation)",
    "Caulking Gun",
    "Glue Gun",
    "Paint Roller",
    "Paint Brush (standard)",

    # Cleaning/household items with obscure measurements
    "Dustpan",
    "Feather Duster",
    "Fly Swatter",
    "Mouse Trap",
    "Fire Alarm",
    "Smoke Detector",
    "Carbon Monoxide Detector",
    "Doorbell",
    "Door Knob",
    "Padlock",
    "Key (house)",
    "Tape (duct tape roll)",

    # Obscure sports/fitness equipment
    "Fencing Sword (epee)",
    "Squash Racket",
    "Table Tennis Paddle",
    "Table Tennis Ball",
    "Badminton Shuttlecock",
    "Water Polo Ball",
    "Handball (team sport)",
    "Golf Tee",
    "Arrow (carbon)",
    "Resistance Band",
    "Foam Roller",
    "Balance Board",
    "Pull-up Bar (doorframe)",
    "Hurdles (track and field)",
    "Shot Put",
    "Discus",
    "Javelin",
    "Hammer (track and field)",
    "Pole Vault Pole",
    "High Jump Bar",
    "Pommel Horse",
    "Vault Table (gymnastics)",
    "Rings (gymnastics)",
    "Uneven Bars (gymnastics)",
    "Gymnastics Balance Beam",
    "Swimming Lane Rope",
    "Baseball Home Plate",
    "Football Endzone",
    "Pogo Stick",
    "Boomerang",

    # Obscure garden/outdoor items
    "Hose Reel",
    "Sprinkler Head",
    "Lawn Aerator",
    "Compost Bin",
    "Bird Feeder",
    "Post Hole Digger",
    "Stump Grinder",
    "Hedge Trimmer",

    # Obscure vehicles/craft
    "Pedicab",
    "Airport Luggage Cart Train",
    "Airport Snow Blower",
    "Airport Fuel Truck",
    "Passenger Airplane Stairs",
    "Aircraft Tug",
    "Autogyro",
    "Microlight Aircraft",
    "Pilot Boat",
    "Boston Whaler (17ft)",
    "Rigid Inflatable Boat (RIB)",
    "Airboat",
    "Cable Laying Ship",
    "Ro-Ro Ship (roll-on/roll-off)",
    "LNG Tanker",
    "Cutter (US Coast Guard)",
    "Dredge Ship",
    "Research Vessel (medium)",
    "Canoe (whitewater)",
    "Kite Surfboard",
    "Chinese Junk Sailboat",
    "Manure Spreader",
    "Seed Drill",
    "Hay Baler",
    "Base Jump Wingsuit",

    # Obscure industrial/equipment
    "Pipeline Pig (inspection)",
    "Cotton Gin (industrial)",
    "Paper Making Machine",
    "Assembly Line (100m section)",
    "Conveyor Belt (warehouse)",
    "Car Crusher",
    "Soybean Silo",
    "Oil Pump Jack (pumpjack)",
    "Derrick (oil drilling)",
    "Irrigation Pivot (1/4 mile)",
    "ATM Network Server",
    "Particle Accelerator (desktop)",
    "Lathe Machine (large)",
    "CNC Milling Machine",
    "Industrial Boiler",
    "Transformers (power grid)",
    "Hydroelectric Generator",
    "Electric Motor (industrial)",

    # Obscure structure items (too generic or abstract as measurements)
    "Sandbox",
    "Fence (wood, 6ft, per 10m)",
    "Chain-link Fence (100m section)",
    "Patio (average residential)",
    "Deck (residential, average)",
    "Closet (walk-in, average)",
    "Elevator Shaft (standard)",
    "Staircase (standard flight)",
    "Guard Rail (100m section)",
    "Speed Bump",
    "Bollard",

    # Obscure building/venue items
    "Ceramic Tile (12x12)",
    "Glass Pane (standard window)",
    "Roll of Carpet (12ft)",
    "Plywood Sheet",

    # Obscure animals (very niche)
    "Thorny Devil",
    "Nile Monitor",
    "Gharial",
    "Freshwater Crocodile",
    "Cape Fur Seal",
    "Weddell Seal",
    "Leopard Seal",
    "Fur Seal",
    "Ringed Seal",
    "Harp Seal",
    "Walking Stick Insect",
    "Atlas Moth",
    "Giant African Millipede",
    "Hercules Beetle",
    "Nautilus",

    # Obscure structures/places
    "Piazza Navona Rome",
    "Olympic Park Munich",
    "Bois de Boulogne Paris",
    "Canary Wharf Tower (London)",
    "CCTV Headquarters (Beijing)",
    "Bank of China Tower (HK)",
    "Obelisk of Buenos Aires",
    "Cleopatra's Needle (London)",
    "Colossus of Rhodes (estimated)",
    "Excalibur Sword (typical reproduction)",
    "Stonehenge Altar Stone",
    "Stonehenge Sarsen Stone",

    # Misc obscure
    "Shopping Bag (paper, full)",
    "Cardboard Box (medium)",
    "Cardboard Box (large moving)",
    "Pallet (wooden)",
    "Film Reel (35mm)",
    "Cassette Tape",
    "VHS Tape",
    "CD",
    "Game Boy",
    "Vinyl Record (LP)",
    "Inline Skates (pair)",
    "Ice Skates (pair)",
    "Snowshoe (pair)",
    "Ski Pole",
    "Automatic Sliding Door",
    "Revolving Door",
    "Garage Door (standard 2-car)",
    "Lighthouse Lens (Fresnel, first order)",
    "Portable Restroom Trailer",
    "Parking Sign",
    "Newspaper (full Sunday edition)",
    "Mailman Bag (full)",
    "Pool Noodle",
    "Kayak Paddle",
    "Canoe Paddle",
    "Wheelie Bin (120L)",
    "Cash Register",
    "Safe (home, small)",
    "Geiger Counter",
    "Seismograph",
    "Satellite Dish (backyard)",
    "Satellite Dish (large telecom)",
    "Radar Antenna (weather)",
    "Industrial Robot Arm",
    "Oscilloscope",
    "Baby Bathtub",
    "Playpen",
}

# ── Items to RENAME in world_measurements.json ──────────────────────────────

MEASUREMENTS_RENAME = {
    "Oil Barrel (55 gallon)": "Oil Barrel",
    "Bucket (standard 5 gallon)": "Bucket",
    "Rain Barrel (55 gallon)": "Rain Barrel",
    "Wooden Ladder (6ft)": "Wooden Ladder",
    "Ladder (extension, 24ft)": "Extension Ladder",
    "Tape Measure (25ft)": "Tape Measure",
    "Kettlebell (24kg)": "Kettlebell",
    "Trampoline (14ft)": "Trampoline",
    "Pool Table (9ft)": "Pool Table",
    "Large Hadron Collider (circumference)": "Large Hadron Collider",
    "Pipeline (oil, 1 mile section)": "Oil Pipeline (1 mile section)",
    "Escalator (standard)": "Escalator",
}

# ── Items to REMOVE from world_durations.json ───────────────────────────────

DURATIONS_REMOVE = {
    # SI/standard time units (the item IS a time unit)
    "Planck time",
    "Attosecond (shortest laser pulse)",
    "Femtosecond",
    "Picosecond",
    "Nanosecond",
    "Microsecond",
    "Millisecond",
    "One second",
    "One minute",
    "One hour",
    "One day (solar)",
    "One week",
    "One month (average)",
    "One year (Julian)",
    "One decade",
    "One century",
    "One millennium",

    # Duplicates
    "Usain Bolt's 100m world record",      # dup of "100m world record (men)"
    "Moon to Earth radio signal delay",     # dup of "Moonlight travel to Earth"
    "Hiroshima bomb detonation (fission)",  # dup of "Nuclear fission chain reaction"
    "REM sleep cycle",                       # overlaps with "Human sleep cycle"
    "Pompeii burial (79 AD eruption)",       # dup of "Pompeii eruption (Vesuvius 79 AD)"
    "Voyager 1 travel to heliopause",        # dup of "Voyager 1 reaching interstellar space"

    # Not commonly known / obscure
    "Nuclear fission chain reaction",
    "Photosynthesis (light reactions)",
    "Speed of fastest human nerve signal",
    "Time light takes to cross human hair width",
    "Minimum human reaction to pain",
    "Nuclear bomb assembly time (modern)",
    "Milankovitch cycle (eccentricity)",
    "Precession of Earth's axis",
    "Great Oxidation Event duration",
    "Snowball Earth glaciation",
    "Average time to fall in love",
    "Time for Earth to rotate 1 degree",
    "Radioactive decay (Carbon-14 half-life)",
    "Radioactive decay (Uranium-238 half-life)",
    "Gamma-ray burst (short)",
    "Gamma-ray burst (long)",
    "Pulsar rotation period",
    "Geomagnetic storm",
    "Solar flare (X-class)",
    "Speed of nerve impulse across body",
    "Cambrian explosion duration",
    "Permian mass extinction",
}


def process_file(path, remove_set, rename_map):
    with open(path) as f:
        items = json.load(f)

    before = len(items)
    removed = []
    renamed = []

    result = []
    for item in items:
        name = item["name"]
        if name in remove_set:
            removed.append(name)
            continue
        if name in rename_map:
            renamed.append((name, rename_map[name]))
            item["name"] = rename_map[name]
        result.append(item)

    after = len(result)

    with open(path, "w") as f:
        json.dump(result, f, indent=2, ensure_ascii=False)
        f.write("\n")

    return before, after, removed, renamed


def main():
    m_path = DATA_DIR / "world_measurements.json"
    d_path = DATA_DIR / "world_durations.json"

    print("Processing world_measurements.json...")
    b, a, removed, renamed = process_file(m_path, MEASUREMENTS_REMOVE, MEASUREMENTS_RENAME)
    print(f"  {b} → {a} items ({b - a} removed)")
    if renamed:
        print(f"  Renamed {len(renamed)} items:")
        for old, new in renamed:
            print(f"    {old} → {new}")
    if removed:
        print(f"  Removed {len(removed)} items")

    # Check for items in remove set that weren't found
    found_names = {r for r in removed}
    missing = MEASUREMENTS_REMOVE - found_names
    if missing:
        print(f"\n  WARNING: {len(missing)} items in remove list not found:")
        for name in sorted(missing):
            print(f"    - {name}")

    print(f"\nProcessing world_durations.json...")
    b, a, removed, renamed = process_file(d_path, DURATIONS_REMOVE, {})
    print(f"  {b} → {a} items ({b - a} removed)")
    if removed:
        print(f"  Removed {len(removed)} items")

    missing = DURATIONS_REMOVE - set(removed)
    if missing:
        print(f"\n  WARNING: {len(missing)} items in remove list not found:")
        for name in sorted(missing):
            print(f"    - {name}")


if __name__ == "__main__":
    main()
