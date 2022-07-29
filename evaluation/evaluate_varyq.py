import numpy as np
import matplotlib.pyplot as plt
import scipy.stats as st
import seaborn as sns
import pandas as pd

# sns.set_theme(style="darkgrid")

folder = "data/"

study = "q,k="
z = [2, 4, 8]
xlims = [0, 1]
xlabel = "Adversary proportion"

printAnalytical = True
ylimsTips = [0, 350]
ylimsOrphan = [1e-7, 1]
cutdata = .5

# Colors
BG_WHITE = "#fbf9f4"
GREY_LIGHT = "#b4aea9"
GREY50 = "#7F7F7F"
BLUE_DARK = "#1B2838"
BLUE = "#2a475e"
BLACK = "#282724"
GREY_DARK = "#747473"
RED_DARK = "#850e00"
# Colors taken from Dark2 palette in RColorBrewer R library
COLOR_SCALE = ["#1B9E77", "#D95F02", "#7570B3"]


def main():
    evaluate1(1)  # tips
    # evaluate1(2)  # orphanage


def evaluate1(analysisType):

    fig, ax = plt.subplots()

    for j in np.arange(len(z)):
        folderdata = "../data/"+study+str(z[j])+",lam=100/"
        if analysisType == 1:
            print("------------ Tips -----------")
            filenamedata = folderdata+"tips_"
            ylabel = "Number of tips"
            fileSaveFig = folder+'tips.png'
        elif analysisType == 2:
            print("------------ Orphans -----------")
            filenamedata = folderdata+"orphantips_"
            ylabel = "Orphanage rate"
            fileSaveFig = folder+'orphanage.png'

        X = loadColumn(folderdata+"params", 0, 0)
        print("+++++++++++++ "+str(z[j])+" +++++++++++++++")
        print("X= ", X)
        print("Length of X="+str(len(X)))
        y = X*0.
        yQ1 = X*0.
        yQ3 = X*0.
        yMin = X*0.
        yMax = X*0.

        for i in np.arange(len(X)):
            y_data = loadColumn(filenamedata+str(i), 1, 2)
            y_data = y_data[int(len(y_data)*cutdata):]
            y[i] = np.mean(y_data)
            samples = 10
            for m in np.arange(samples):
                # use sample to calculate yQ1 and yQ3
                y_data = loadColumn(filenamedata+str(i), 4+m, 2)
                y_data = y_data[int(len(y_data)*cutdata):]
                df = pd.DataFrame(y_data, columns=['data'])
                dfStats = df['data'].describe()
                yQ1[i] += dfStats['25%']/float(samples)
                yQ3[i] += dfStats['75%']/float(samples)
        # 0 values make no sense
        y[y == 0] = np.nan
        yQ3[yQ1 == 0] = np.nan
        yQ1[yQ1 == 0] = np.nan
        yMax[yMin == 0] = np.nan
        yMin[yMin == 0] = np.nan
        sns.lineplot(x=X, y=y, label="k="+str(z[j]))
        Label = ""
        if j == len(z)-1:  # only if last last one
            Label = "25% to 75%\nquantiles"
        plt.fill_between(X, yQ1, yQ3, color='b',
                         alpha=0.2, label=Label)
        # plt.fill_between(X, yMin, yQ1, color='r',
        #                  alpha=0.1, label="Min to Max")
        plt.fill_between(X, yQ3, yMax, color='r',
                         alpha=0.1)
        if printAnalytical & analysisType == 1:
            xL,  Lsimple = getAnalyticalCurve(folderdata, z[j])
            print("xL=", xL)
            Label = ""
            if j == len(z)-1:  # only if last last one
                Label = "Analytical"
            sns.lineplot(x=xL, y=Lsimple, color="red",
                         label=Label, linestyle="dashed")

    # general stuff
    plt.ylim(ylimsTips)
    if analysisType == 2:
        plt.yscale('log')
        plt.ylim(ylimsOrphan)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.xlim(xlims)
    plt.legend()
    sns.despine()
    plt.grid()
    plt.savefig(fileSaveFig, format='png')
    plt.clf()


def loadColumn(filename, column, skiprows):
    try:
        filestr = filename+".csv"
        f = open(filestr, "r")
        data = np.loadtxt(f, delimiter=";",
                          skiprows=skiprows, usecols=(column))
        return data
    except FileNotFoundError:
        print(filestr)
        print("File not found.")
        return []


def getAnalyticalCurve(folder, k):
    lam = 100.
    h = 1.

    # load X data
    Xdata = loadColumn(folder+"params", 0, 0)*1.1
    X = np.arange(50)/50.
    X = max(Xdata)*np.ones(len(X))-(max(Xdata)-min(Xdata))*X
    q = X
    p = 1-q

    # simple equation instead
    Lsimple = p*k/(p*k-1)*h*lam*np.ones(len(X))
    Lsimple[Lsimple < 0] = np.NaN

    return X,  Lsimple


# needs to be at the very end of the file
if __name__ == '__main__':
    main()
