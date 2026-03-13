(function() {
    const getBasePath = function() {
        const path = window.location.pathname;
        // 移除末尾的斜杠，但保留开头和中间的
        let basePath = path.endsWith('/') ? path.slice(0, -1) : path;
        // 确保以 / 开头
        if (!basePath.startsWith('/')) {
            basePath = '/' + basePath;
        }
        return basePath;
    };

    const basePath = getBasePath();
    const api = function(endpoint) {
        return basePath + endpoint;
    };

    let currentPath = '/';
    let selectedFiles = [];

    const loginPage = document.getElementById('login-page');
    const filePage = document.getElementById('file-page');
    const loginForm = document.getElementById('login-form');
    const fileList = document.getElementById('file-tbody');
    const breadcrumb = document.getElementById('breadcrumb');
    const fileInput = document.getElementById('file-input');
    const dropZone = document.getElementById('drop-zone');

    let editorInstance = null;
    let currentEditPath = '';

    function getModeFromPath(path) {
        if (!path) return 'text/plain';
        const ext = path.split('.').pop().toLowerCase();
        const modeMap = {
            'js': 'javascript',
            'jsx': 'javascript',
            'ts': 'javascript',
            'tsx': 'javascript',
            'json': 'javascript',
            'css': 'css',
            'scss': 'css',
            'less': 'css',
            'html': 'htmlmixed',
            'htm': 'htmlmixed',
            'xml': 'xml',
            'py': 'python',
            'md': 'markdown',
            'markdown': 'markdown',
            'go': 'go',
            'java': 'text/plain',
            'c': 'text/plain',
            'cpp': 'text/plain',
            'h': 'text/plain',
            'sh': 'shell',
            'bash': 'shell',
            'yaml': 'yaml',
            'yml': 'yaml',
            'sql': 'sql',
            'php': 'php',
            'rb': 'ruby',
            'rs': 'rust',
            'vue': 'htmlmixed',
            'svelte': 'htmlmixed'
        };
        return modeMap[ext] || 'text/plain';
    }

    let editorLoaded = false;
    let editorLoading = false;

    function loadEditor() {
        return new Promise((resolve, reject) => {
            if (editorLoaded) {
                resolve();
                return;
            }
            if (editorLoading) {
                const checkLoaded = setInterval(() => {
                    if (editorLoaded) {
                        clearInterval(checkLoaded);
                        resolve();
                    }
                }, 100);
                return;
            }
            
            editorLoading = true;
            const loader = document.getElementById('editor-loader');
            loader.classList.remove('hidden');
            
            const cssUrl = 'https://cdn.jsdelivr.net/npm/codemirror@5/lib/codemirror.min.css';
            const jsUrl = 'https://cdn.jsdelivr.net/npm/codemirror@5/lib/codemirror.min.js';
            
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.href = cssUrl;
            document.head.appendChild(link);
            
            const script = document.createElement('script');
            script.src = jsUrl;
            script.onload = () => {
                editorLoaded = true;
                editorLoading = false;
                loader.classList.add('hidden');
                resolve();
            };
            script.onerror = () => {
                editorLoading = false;
                loader.classList.add('hidden');
                reject(new Error('Failed to load editor'));
            };
            document.head.appendChild(script);
        });
    }

    function initEditor(content, path, readonly) {
        const textarea = document.getElementById('editor');
        
        if (editorInstance) {
            editorInstance.toTextArea();
            editorInstance = null;
        }
        
        textarea.style.display = 'block';
        
        const mode = getModeFromPath(path);
        
        editorInstance = CodeMirror.fromTextArea(textarea, {
            mode: mode,
            theme: 'default',
            lineNumbers: true,
            lineWrapping: true,
            indentWithTabs: true,
            indentUnit: 4,
            tabSize: 4,
            readOnly: readonly,
            autofocus: true,
            extraKeys: {
                'Ctrl-S': function() {
                    if (!readonly && currentEditPath) {
                        saveFile(currentEditPath);
                    }
                },
                'Cmd-S': function() {
                    if (!readonly && currentEditPath) {
                        saveFile(currentEditPath);
                    }
                }
            }
        });
        
        editorInstance.setValue(content || '');
    }

    function getEditorContent() {
        return editorInstance ? editorInstance.getValue() : '';
    }

    function init() {
        checkAuth();
        bindEvents();
    }

    function checkAuth() {
        fetch(api('/api/check'))
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    showFilePage();
                    loadFiles(currentPath);
                } else {
                    showLoginPage();
                }
            })
            .catch(() => showLoginPage());
    }

    function showLoginPage() {
        loginPage.classList.remove('hidden');
        filePage.classList.add('hidden');
    }

    function showFilePage() {
        loginPage.classList.add('hidden');
        filePage.classList.remove('hidden');
    }

    function bindEvents() {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const password = document.getElementById('password').value;
            
            fetch(api('/api/login'), {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({password: password})
            })
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    showFilePage();
                    loadFiles(currentPath);
                } else {
                    alert(data.msg || '登录失败');
                }
            });
        });

        document.getElementById('btn-upload').addEventListener('click', () => fileInput.click());
        document.getElementById('btn-new-folder').addEventListener('click', showNewFolderModal);
        document.getElementById('btn-refresh').addEventListener('click', () => loadFiles(currentPath));
        document.getElementById('btn-logout').addEventListener('click', logout);

        fileInput.addEventListener('change', handleUpload);

        document.addEventListener('dragover', e => {
            e.preventDefault();
            dropZone.classList.add('active');
        });
        document.addEventListener('dragleave', e => {
            if (e.target === dropZone) {
                dropZone.classList.remove('active');
            }
        });
        document.addEventListener('drop', e => {
            e.preventDefault();
            dropZone.classList.remove('active');
            if (e.dataTransfer.files.length > 0) {
                uploadFiles(e.dataTransfer.files);
            }
        });

        document.getElementById('check-all').addEventListener('change', function() {
            const checks = document.querySelectorAll('.file-check');
            checks.forEach(c => c.checked = this.checked);
            updateSelectedFiles();
        });
    }

    function loadFiles(path) {
        currentPath = path;
        fetch(api('/api/list?path=' + encodeURIComponent(path)))
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    renderFileList(data.data);
                    renderBreadcrumb(path);
                }
            });
    }

    function renderFileList(files) {
        fileList.innerHTML = '';
        
        files.forEach(file => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td><input type="checkbox" class="file-check" data-path="${file.path}"></td>
                <td>
                    <span class="file-icon">${file.isDir ? '📁' : '📄'}</span>
                    <span class="file-name" data-path="${file.path}" data-isdir="${file.isDir}">${file.name}</span>
                </td>
                <td>${file.isDir ? '-' : formatSize(file.size)}</td>
                <td>${file.modTime}</td>
                <td class="actions">
                    ${!file.isDir ? '<button class="btn-view">查看</button><button class="btn-edit">编辑</button>' : ''}
                    <button class="btn-download">下载</button>
                    <button class="btn-rename">重命名</button>
                    <button class="btn-copy">复制</button>
                    <button class="btn-move">移动</button>
                    <button class="btn-delete">删除</button>
                </td>
            `;
            
            const nameEl = tr.querySelector('.file-name');
            if (file.isDir) {
                nameEl.addEventListener('click', () => loadFiles(file.path));
            } else {
                nameEl.addEventListener('click', () => viewFile(file.path));
            }

            const actions = tr.querySelector('.actions');
            if (!file.isDir) {
                actions.querySelector('.btn-view')?.addEventListener('click', () => viewFile(file.path));
                actions.querySelector('.btn-edit')?.addEventListener('click', () => editFile(file.path));
                actions.querySelector('.btn-download')?.addEventListener('click', () => downloadFile(file.path));
            }
            actions.querySelector('.btn-rename').addEventListener('click', () => showRenameModal(file.path, file.name));
            actions.querySelector('.btn-copy').addEventListener('click', () => showCopyModal(file.path));
            actions.querySelector('.btn-move')?.addEventListener('click', () => showMoveModal(file.path));
            actions.querySelector('.btn-delete').addEventListener('click', () => deleteFile(file.path));

            fileList.appendChild(tr);
        });
    }

    function renderBreadcrumb(path) {
        breadcrumb.innerHTML = '';
        const parts = path.split('/').filter(p => p);
        
        const rootSpan = document.createElement('span');
        rootSpan.className = 'path-item';
        rootSpan.textContent = '根目录';
        rootSpan.addEventListener('click', () => loadFiles('/'));
        breadcrumb.appendChild(rootSpan);

        let accPath = '';
        parts.forEach(part => {
            accPath += '/' + part;
            const sep = document.createElement('span');
            sep.textContent = ' > ';
            breadcrumb.appendChild(sep);

            const span = document.createElement('span');
            span.className = 'path-item';
            span.textContent = part;
            const clickPath = accPath;  // Capture current value
            span.addEventListener('click', () => loadFiles(clickPath));
            breadcrumb.appendChild(span);
        });
    }

    function formatSize(bytes) {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
        return (bytes / 1024 / 1024).toFixed(1) + ' MB';
    }

    function viewFile(path) {
        loadEditor().then(() => {
            fetch(api('/api/read?path=' + encodeURIComponent(path)))
                .then(r => r.json())
                .then(data => {
                    if (data.code === 0) {
                        const modal = document.getElementById('edit-modal');
                        initEditor(data.data, path, true);
                        document.getElementById('btn-save-edit').classList.add('hidden');
                        modal.classList.remove('hidden');
                    }
                });
        }).catch(() => {
            alert('编辑器加载失败');
        });
    }

    function editFile(path) {
        loadEditor().then(() => {
            fetch(api('/api/read?path=' + encodeURIComponent(path)))
                .then(r => r.json())
                .then(data => {
                    if (data.code === 0) {
                        const modal = document.getElementById('edit-modal');
                        currentEditPath = path;
                        initEditor(data.data, path, false);
                        document.getElementById('btn-save-edit').classList.remove('hidden');
                        document.getElementById('btn-save-edit').onclick = () => saveFile(path);
                        modal.classList.remove('hidden');
                    }
                });
        }).catch(() => {
            alert('编辑器加载失败');
        });
    }

    function saveFile(path) {
        const content = getEditorContent();
        fetch(api('/api/write'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({path: path, content: content})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                document.getElementById('edit-modal').classList.add('hidden');
                if (editorInstance) {
                    editorInstance.toTextArea();
                    editorInstance = null;
                }
                currentEditPath = '';
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    function downloadFile(path) {
        window.location.href = api('/api/download?path=' + encodeURIComponent(path));
    }

    function handleUpload(e) {
        uploadFiles(e.target.files);
    }

    function uploadFiles(files) {
        const formData = new FormData();
        formData.append('path', currentPath);
        
        for (let i = 0; i < files.length; i++) {
            formData.append('file', files[i]);
        }

        fetch(api('/api/upload'), {
            method: 'POST',
            body: formData
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    function showNewFolderModal() {
        document.getElementById('new-folder-modal').classList.remove('hidden');
    }

    document.getElementById('btn-create-folder').addEventListener('click', function() {
        const name = document.getElementById('new-folder-name').value;
        if (!name) return;

        fetch(api('/api/create'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({path: currentPath, name: name})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                document.getElementById('new-folder-modal').classList.add('hidden');
                document.getElementById('new-folder-name').value = '';
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    });

    document.getElementById('btn-cancel-folder').addEventListener('click', () => {
        document.getElementById('new-folder-modal').classList.add('hidden');
    });

    let renamePath = '';
    function showRenameModal(path, name) {
        renamePath = path;
        document.getElementById('rename-input').value = name;
        document.getElementById('rename-modal').classList.remove('hidden');
    }

    document.getElementById('btn-rename').addEventListener('click', function() {
        const name = document.getElementById('rename-input').value;
        if (!name) return;

        fetch(api('/api/rename'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({path: renamePath, name: name})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                document.getElementById('rename-modal').classList.add('hidden');
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    });

    document.getElementById('btn-cancel-rename').addEventListener('click', () => {
        document.getElementById('rename-modal').classList.add('hidden');
    });

    function deleteFile(path) {
        if (!confirm('确定要删除吗？')) return;

        fetch(api('/api/delete'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({paths: [path]})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    function showCopyModal(src) {
        const dst = prompt('请输入目标路径（如 /folder/name）:');
        if (!dst) return;

        fetch(api('/api/copy'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({src: src, dst: dst})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    function showMoveModal(src) {
        const dst = prompt('请输入目标路径（如 /folder/name）:');
        if (!dst) return;

        fetch(api('/api/move'), {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({src: src, dst: dst})
        })
        .then(r => r.json())
        .then(data => {
            if (data.code === 0) {
                loadFiles(currentPath);
            } else {
                alert(data.msg);
            }
        });
    }

    document.getElementById('btn-cancel-edit').addEventListener('click', () => {
        document.getElementById('edit-modal').classList.add('hidden');
        if (editorInstance) {
            editorInstance.toTextArea();
            editorInstance = null;
        }
        currentEditPath = '';
    });

    function logout() {
        fetch(api('/api/logout'), {method: 'POST'})
            .then(() => {
                showLoginPage();
            });
    }

    function updateSelectedFiles() {
        selectedFiles = Array.from(document.querySelectorAll('.file-check:checked')).map(c => c.dataset.path);
    }

    init();
})();
