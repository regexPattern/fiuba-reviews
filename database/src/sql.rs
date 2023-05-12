pub trait Sql {
    fn sanitizar(&self) -> String;
}

impl Sql for &str {
    fn sanitizar(&self) -> String {
        self.replace('\'', "''")
    }
}

impl Sql for String {
    fn sanitizar(&self) -> Self {
        self.as_str().sanitizar()
    }
}
