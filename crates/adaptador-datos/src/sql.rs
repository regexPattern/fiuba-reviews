pub trait Sql {
    fn sanitizar_sql(&self) -> String;
}

impl Sql for &str {
    fn sanitizar_sql(&self) -> String {
        self.replace('\'', "''")
    }
}

impl Sql for String {
    fn sanitizar_sql(&self) -> Self {
        self.as_str().sanitizar_sql()
    }
}
