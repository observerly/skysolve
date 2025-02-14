-- Create "stars" table
CREATE TABLE `stars` (`id` text NOT NULL, `designation` text NOT NULL, `x` float NOT NULL, `y` float NOT NULL, `ra` float NOT NULL, `dec` float NOT NULL, `intensity` float NOT NULL, `pixel` integer NOT NULL, PRIMARY KEY (`designation`));
-- Create index "idx_pixel" to table: "stars"
CREATE INDEX `idx_pixel` ON `stars` (`pixel`);
