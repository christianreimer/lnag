# Large Number Analogy Generator (lnag)

Large Number Analogy Generator is a service to help visualize large or small numbers by comparing them to physical concepts. For example, Apple has sold over 3 million iPhones. If they were all stacked on top of each other then the wobbly tower of phones would reach more than half way to the moon!

The service can express numbers in terms of duration, length, height, weight, and volume. The caller can indicate which dimension they want to use (or one will be picked randomly) and if they are looking to express how small or large something is.

The service has a library of concepts that can be matched to produce the visualization. For example, a list of how think items are and a list of distances, can produce "N <items> placed next to each other would reach from <start> to <end>"

The service is exposed as a simple HTTP endpoint and produces a JSON response.
