package apex

object TypeSafeEquality {
  implicit class Equal[A](left: A) {
    def =?(right: A): Boolean = left == right
    def !=?(right:A): Boolean = left != right
  }
}