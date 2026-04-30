# Third-Party Open Source Notices — Anybackup CLI

This notice enumerates the third-party open source components introduced by the Go module under `repos/Anybackup/CLI`, both as direct and transitive dependencies compiled into the distributed binary, together with their respective license information. The CLI source code itself is distributed under SSPL-1.0 together with the root [LICENSE](../LICENSE) of this repository, and is therefore out of scope for this file.

## 1. Scope

In scope:

- Go modules declared in [`go.mod`](./go.mod) that are linked into the CLI binary at build time, including both direct and indirect dependencies.
- Any additional third-party libraries statically embedded in artifacts produced by `go build` for this CLI.

Out of scope (see the repository knowledge entry *"Pure protocol-level invocation does not trigger third-party open source declaration obligations"*):

- The CLI talks to the Anybackup console over HTTPS. This is a runtime protocol interop and does not constitute a dependency on the server's source code or binaries; those components are declared in their own charts or repositories.

## 2. Third-Party Component Inventory

| Component | Version | License | How It Is Introduced | Location in Module | Upstream Project |
|---|---|---|---|---|---|
| github.com/spf13/cobra | v1.9.1 | Apache License 2.0 | Direct Go module dependency | [go.mod](./go.mod) / [go.sum](./go.sum) | https://github.com/spf13/cobra |
| github.com/spf13/pflag | v1.0.6 | BSD 3-Clause License | Indirect Go module dependency (required by cobra) | [go.mod](./go.mod) / [go.sum](./go.sum) | https://github.com/spf13/pflag |
| github.com/inconshreveable/mousetrap | v1.1.0 | Apache License 2.0 | Indirect Go module dependency (required by cobra on Windows) | [go.mod](./go.mod) / [go.sum](./go.sum) | https://github.com/inconshreveable/mousetrap |

## 3. License Texts and References

### 3.1 cobra, mousetrap — Apache License 2.0

Copyright © the cobra authors; copyright © Alan Shreve (mousetrap).

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Canonical license texts:

- https://github.com/spf13/cobra/blob/main/LICENSE.txt
- https://github.com/inconshreveable/mousetrap/blob/master/LICENSE

### 3.2 pflag — BSD 3-Clause License

Copyright (c) 2012 Alex Ogier. All rights reserved.
Copyright (c) 2012 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.

Canonical license text: https://github.com/spf13/pflag/blob/master/LICENSE

## 4. Compliance and Redistribution

- The CLI is NOT a derivative work of the above components. They are linked as Go modules according to their own license terms.
- Any redistribution of the CLI binary, or of container images and release bundles that embed it, MUST carry this notice together with the full Apache License 2.0 and BSD 3-Clause license texts referenced above.
- Introducing a new direct or indirect Go module dependency into the CLI build MUST be accompanied by an update to this file in the same change.

## 5. Maintenance Rules

- When any component listed above is upgraded, the "Version" column in this file MUST be updated in the same change.
- When this module introduces a new third-party dependency (direct or indirect) that is statically linked into the binary, the component, version, license and upstream link MUST be registered in this file.
- Keep the split consistent with the chart-side notices under [`deploy/helm/**/THIRD_PARTY_NOTICES.md`](../deploy/helm): dependencies of the CLI live here; dependencies of a specific chart or service image live next to that chart or image.
- This file is an open source compliance notice maintained inside the code sub-repository and follows the rules stated in `repos/AGENTS.md` and `repos/_rules/`.
