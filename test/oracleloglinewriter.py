from sqlalchemy import create_engine, Table, Column, Integer, Text, Boolean, MetaData, insert, select, Sequence, text
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
db_type = "oracle"  # Change this to "mysql" if you're using MySQL

# Database connection URL
db_url = get_db_url(db_type)

# Initialize SQLAlchemy engine
engine = create_engine(db_url)
metadata = MetaData()

# Define the logs table structure with a sequence for Oracle
logs_table = Table(
    "logs", metadata,
    Column("id", Integer, Sequence('logs_seq', start=1, increment=1), primary_key=True),  # Use sequence for Oracle ID autoincrement
    Column("message", Text, nullable=False),
    Column("processed", Boolean, nullable=False, default=False),
    Column("addcolmn", Text, nullable=False),
)

def initialize_sequence():
    """Creates the sequence for Oracle if it doesn't already exist."""
    if db_type == "oracle":
        try:
            with engine.connect() as connection:
                connection.execute(text("CREATE SEQUENCE logs_seq START WITH 1 INCREMENT BY 1"))
                print("Sequence 'logs_seq' created successfully.")
        except SQLAlchemyError as err:
            # Ignore the error if sequence already exists (ORA-00955)
            if "ORA-00955" in str(err):
                print("Sequence 'logs_seq' already exists.")
            else:
                print(f"Error creating sequence: {err}")

def initialize_database():
    try:
        # Create sequence for Oracle
        initialize_sequence()
        
        # Create the logs table if it doesn't exist
        metadata.create_all(engine)
        print("Table 'logs' checked/created successfully.")

        # Insert a log if no rows exist
        with engine.connect() as connection:
            result = connection.execute(select(logs_table.c.id).limit(1))  # Check if there's at least one log
            if result.fetchone() is None:
                ins = insert(logs_table).values(message="This is an example log message", addcolmn="extracolnm", processed=False)
                connection.execute(ins)
                connection.commit()
                print("Inserted example log message into 'logs' table.")
            else:
                print("The 'logs' table already has entries; no new log was inserted.")

    except SQLAlchemyError as err:
        print(f"Database error: {err}")

def writeline2db():
    try:
        # Insert a log message into the table
        with engine.connect() as connection:
            ins = insert(logs_table).values(message="This is an example log message", addcolmn="extracolnm", processed=False)
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
        
        if db_type == "oracle":
            # Try to drop the sequence, ignoring errors if it doesn't exist
            with engine.connect() as connection:
                try:
                    connection.execute(text("DROP SEQUENCE logs_seq"))
                    print("Sequence 'logs_seq' has been dropped successfully.")
                except SQLAlchemyError as err:
                    # Only handle the error if the sequence does not exist
                    if "ORA-02289" in str(err):
                        print("Sequence 'logs_seq' does not exist, so it was not dropped.")
                    else:
                        raise  # Raise other errors

    except SQLAlchemyError as err:
        print(f"Error while dropping the table or sequence: {err}")


if __name__ == "__main__":
    initialize_database()
    writeline2db()
    #delete_table()
