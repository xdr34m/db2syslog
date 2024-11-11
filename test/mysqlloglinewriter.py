from sqlalchemy import create_engine, Table, Column, Integer, Text, Boolean, MetaData, insert, select, Sequence
from sqlalchemy.exc import SQLAlchemyError

def get_db_url(db_type: str):
    """Returns the database URL based on the db_type ('mysql' or 'oracle')."""
    if db_type == "mysql":
        # MySQL connection string
        return "mysql+pymysql://user:password@127.0.0.1:3306/logsdb"
    elif db_type == "oracle":
        # Oracle connection string, replace with your own details
        return "oracle+oracledb://SYSTEM:password@127.0.0.1:1521/?service_name=logsdb"
    else:
        raise ValueError("Unsupported database type. Please choose 'mysql' or 'oracle'.")

# Specify the database type (either 'mysql' or 'oracle')
db_type = "mysql"  # You can change this to "mysql" if you're using MySQL

# Database connection URL
db_url = get_db_url(db_type)

# Initialize SQLAlchemy engine
engine = create_engine(db_url)
metadata = MetaData()

# Define the logs table structure
logs_table = Table(
    "logs", metadata,
    Column("id", Integer, primary_key=True, autoincrement=True),
    Column("message", Text, nullable=False),
    Column("processed", Boolean, nullable=False, default=False),
    Column("addcolmn", Text, nullable=False),
)

def get_db_url(db_type: str):
    """Returns the database URL based on the db_type ('mysql' or 'oracle')."""
    if db_type == "mysql":
        # MySQL connection string
        return "mysql+pymysql://user:password@127.0.0.1:3306/logsdb"
    elif db_type == "oracle":
        # Oracle connection string, replace with your own details
        return "oracle+python-oracledb://SYSTEM:password@127.0.0.1:1521/?service_name=logsdb"
    else:
        raise ValueError("Unsupported database type. Please choose 'mysql' or 'oracle'.")

def initialize_database():
    try:
        # Create the logs table if it doesn't exist
        metadata.create_all(engine)
        print("Table 'logs' checked/created successfully.")

        # Check if there are any logs in the table
        with engine.connect() as connection:
            # Corrected query to select columns
            result = connection.execute(select(logs_table.c.id).limit(1))  # Select 'id' column
            if result.fetchone() is None:
                # Insert an example log message if the table is empty
                ins = insert(logs_table).values(message="This is an example log message",addcolmn="extracolnm", processed=False)
                connection.execute(ins)
                connection.commit()
                print("Inserted example log message into 'logs' table.")
            else:
                print("The 'logs' table already has entries; no new log was inserted.")

    except SQLAlchemyError as err:
        print(f"Database error: {err}")

def writeline2db():
    try:
        # Create the logs table if it doesn't exist
        metadata.create_all(engine)
        print("Table 'logs' checked/created successfully.")

        # Check if there are any logs in the table
        with engine.connect() as connection:

            # Insert an example log message if the table is empty
            ins = insert(logs_table).values(message="This is an example log message",addcolmn="extracolnm", processed=False)
            connection.execute(ins)
            connection.commit()
            print("Inserted example log message into 'logs' table.")


    except SQLAlchemyError as err:
        print(f"Database error: {err}")
        
def delete_table():
    try:
        # Drop the 'logs' table
        logs_table.drop(engine)
        print("Table 'logs' has been dropped successfully.")

    except SQLAlchemyError as err:
        print(f"Error while dropping the table: {err}")
    
if __name__ == "__main__":
    initialize_database()
    writeline2db()
    #delete_table()