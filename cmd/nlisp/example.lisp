;(defun fib (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2)))))

(setq fib (lambda (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2))))))
(println (timeit (fib 35)))
(setq fib (compile fib))
(println (timeit (fib 35)))

