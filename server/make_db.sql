/* Usage: sqlite3 main.db < make_db.sql 
 * INSERT INTO tag (name) VALUES ("code-the-change");
 */

CREATE TABLE tag(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(255));

CREATE TABLE project(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title VARCHAR(255),
	github VARCHAR(255),
	organization VARCHAR(255),
	description TEXT);

CREATE TABLE project_tag(
	projectID,
	tagID,
	FOREIGN KEY(projectID) REFERENCES project(id),
	FOREIGN KEY(tagID) REFERENCES tag(id));