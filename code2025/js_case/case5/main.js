// main.js
const startBtn = document.getElementById('startBtn');
const statusEl = document.getElementById('status');
const resultEl = document.getElementById('result');

// 创建工作线程（注意：Worker 脚本路径需符合同源策略）
const dataWorker = new Worker('./dataProcessor.worker.js', { type: 'module' });

// 生成测试数据（10 万条随机数）
function generateTestData() {
  const data = [];
  for (let i = 0; i < 100000; i++) {
    data.push(Math.random() * 1000); // 0-1000 的随机数
  }
  return data;
}

// 点击按钮开始处理
startBtn.addEventListener('click', () => {
  statusEl.textContent = '状态：处理中...（主线程仍可操作）';
  resultEl.textContent = '';
  startBtn.disabled = true;

  // 生成数据并发送给 Worker 处理
  const testData = generateTestData();
  dataWorker.postMessage(testData); // 向 Worker 发送消息
});

// 接收 Worker 处理后的结果
dataWorker.onmessage = (e) => {
  const { sum, avg, max } = e.data;
  statusEl.textContent = '状态：处理完成';
  resultEl.innerHTML = `
    <p>数据总量：100,000 条</p>
    <p>总和：${sum.toFixed(2)}</p>
    <p>平均值：${avg.toFixed(2)}</p>
    <p>最大值：${max.toFixed(2)}</p>
  `;
  startBtn.disabled = false;
};

// 监听 Worker 错误
dataWorker.onerror = (error) => {
  statusEl.textContent = `状态：处理出错 - ${error.message}`;
  startBtn.disabled = false;
  dataWorker.terminate(); // 出错时终止 Worker
};