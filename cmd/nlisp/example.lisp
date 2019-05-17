(defun fib (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2)))))

(setq fib (lambda (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2))))))

(println (format "2 + 2 = %d" (plus 2 2)))

