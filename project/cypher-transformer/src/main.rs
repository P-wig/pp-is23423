use sqlparser::transformer::transform_cypher_to_sql;
use std::env;
use std::process;

fn main() {
    let args: Vec<String> = env::args().collect();
    
    if args.len() != 2 {
        eprintln!("Usage: cypher_transformer <cypher_query>");
        process::exit(1);
    }
    
    let cypher_query = &args[1];
    
    match transform_cypher_to_sql(cypher_query) {
        Ok(sql) => println!("{}", sql),
        Err(e) => {
            eprintln!("Error: {}", e);
            process::exit(1);
        }
    }
}
