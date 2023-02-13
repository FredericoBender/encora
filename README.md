# encora

    This script make a extraction from the https://careers.encora.com/search?category=All , making a web scrapping and getting all the jobs from the website, also putting them in the database.

## **SPIDERS:**

* encoraSpider.go
  * This script is responsible per extract all the jobs from the encora website and insert the information in the database;
* encoraSpiderVisual.go
  * This script is the same as the encoraSpider.go, but runs with a graphic interface.

## **Installing the dependencies:**

You need to install Chromedriver in order to run the scripts, for that Run the following steps in the:

1. `sudo apt-get install unzip`
2. `wget -N http://chromedriver.storage.googleapis.com/2.29/chromedriver_linux64.zip -P /home/encora/Downloads`
3. `unzip /home/encora/Downloads/chromedriver_linux64.zip -d /home/falqon/Downloads`
4. `rm /home/falqon/Downloads/chromedriver_linux64.zip`
5. `sudo chmod +x /home/encora/Downloads/chromedriver`

## **Running the code:**

1. To run the **encoraSpider** you need to give the destiny database info to extract the data, you can make it with the following examples using Flags.

   * Arguments:
     * --dbUsername: is the username used in the database
     * --dbPassword: is the password of the database
     * --dbName : is the database name
     * --dbHost:  is the host of the database
     * --dbPort: is the port of the database

   1. **example A (In the root file of the project run):**` go run ./cmd/spiders/encoraSpider.go --dbUsername postgres --dbPassword 123 --dbName encora --dbHost localhost --dbPort 5432`
      * In the **example A** all jobs from encora will be extracted and inserted in the postgres database named encora.
   2. **example B (In the root file of the project run):** `go run ./cmd/spiders/encoraSpiderVisual.go `
      * In the **example B**, after run the command, a visual interface appears requesting the same data as the example A, after write the informations of the database in the fields, the script will run with a progress bar until the end.
