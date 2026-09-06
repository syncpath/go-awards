## About
This repository presents a document generation pipeline, Typst-based rendering templates, and PDF diff/finder tooling for producing and validating personalized certificates, diplomas, and gratitude letters at scale for events.

The project is completed during the preparation of Salavat I. Yabbarov's work at SPbPU Institute of Computer Science and Cybersecurity (SPbPU ICSC).

## Authors and contributors
The advisor and contributor Vladimir A. Parkhomenko, Senior Lecturer of SPbPU ICSC.
The main contributor Salavat I. Yabbarov, student of SPbPU ICSC.

## License

Warranty The contributors give no warranty for the using of the software.

---

## Benchmarks

A benchmark was performed for the full Typst document generation pipeline:

1. generate Typst source from configuration and input data;
2. write the generated source to a `.typ` file;
3. invoke `typst compile`;
4. produce the resulting PDF document.

Test environment:

- OS: macOS (`darwin`)
- Architecture: ARM64
- CPU: Apple M5
- Benchmark: `BenchmarkGenerateAndCompileTypst`
- Independent runs: 20

Command:

```bash
go test -bench=BenchmarkGenerateAndCompileTypst -benchmem -count=20
```

### Results

- Mean execution time: 55.808 ms/op
- Standard deviation: 1.710 ms
- Standard error: 0.382 ms
- 95% confidence interval: 55.008–56.608 ms/op
- benchstat: 55.64 ms/op ± 2%
- Memory allocations: 352.5 KiB/op
- Allocations: ~2037 allocs/op

System-level measurements:
- Maximum resident set size: 96.1 MiB
- Peak memory footprint: 17.5 MiB
- Page faults: 6683
- Swaps: 0
- Real time: 23.54 s
- User CPU time: 18.75 s
- System CPU time: 10.89 s
- Average CPU usage: ~126%

The benchmark measures the complete pipeline, including the external typst
process. Therefore, B/op and allocs/op primarily reflect allocations made
by the Go process and do not represent the internal allocations of Typst.
