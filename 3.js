// 猜数字游戏
function guessingGame() {
    const secretNumber = Math.floor(Math.random() * 100) + 1;
    let attempts = 0;
    let guess;

    while (guess !== secretNumber) {
        guess = parseInt(prompt("猜一个1-100之间的数字:"));
        attempts++;

        if (isNaN(guess)) {
            alert("请输入有效的数字!");
        } else if (guess < secretNumber) {
            alert("太小了!再试一次。");
        } else if (guess > secretNumber) {
            alert("太大了!再试一次。");
        } else {
            alert(`恭喜你!答案是${secretNumber}，你用了${attempts}次尝试。`);
        }
    }
}

// 启动游戏
guessingGame();