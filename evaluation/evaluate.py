import numpy as np
import matplotlib.pyplot as plt
import scipy.stats as st
import seaborn as sns
import pandas as pd

sns.set_theme(style="darkgrid")

folder = "data/"

# filename = "q"
# xlims = [0, 1]
# xlabel = "Adversary proportion"
filename = "D,q=.25,k=2,lam=100"
xlims = [1, 19]
xlabel = "Expiration time"
# filename = "k,q=.5"
# xlims = [2, 12]
# xlabel = "Number of parents"

printAnalytical = True
ylims = [1e-7, 1]
folderdata = "../data/"+filename+"/"

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
    evaluate2(1)  # tips
    evaluate1(2)  # orphanage
    evaluate2(2)  # orphanage


def evaluate1(analysisType):
    if analysisType == 1:
        filenamedata = folderdata+"tips_"
        ylabel = "Number of tips"
        fileSaveFig = folder+'tips.png'
        column = 1
    elif analysisType == 2:
        filenamedata = folderdata+"orphantips_"
        ylabel = "Orphanage rate"
        fileSaveFig = folder+'orphanage.png'
        column = 2

    X = loadColumn(folderdata+"params", 0, 0)
    print(X)
    print("Length of X="+str(len(X)))
    y = X*0.
    yQ1 = X*0.
    yQ3 = X*0.
    yMin = X*0.
    yMax = X*0.

    fig, ax = plt.subplots()
    for i in np.arange(len(X)):
        y_data = loadColumn(filenamedata+str(i), column, 2)
        y[i] = np.mean(y_data)
        df = pd.DataFrame(y_data, columns=['data'])
        dfStats = df['data'].describe()
        yQ1[i] = dfStats['25%']
        yQ3[i] = dfStats['75%']
        # 0 values make no sense
        y[y == 0] = np.nan
        yQ3[yQ1 == 0] = np.nan
        yQ1[yQ1 == 0] = np.nan
        yMax[yMin == 0] = np.nan
        yMin[yMin == 0] = np.nan
    sns.lineplot(x=X, y=y, label="Median")
    plt.fill_between(X, yQ1, yQ3, color='b',
                     alpha=0.2, label="25% to 75% quantiles")
    plt.fill_between(X, yMin, yQ1, color='r',
                     alpha=0.1, label="Min to Max")
    plt.fill_between(X, yQ3, yMax, color='r',
                     alpha=0.1)
    if printAnalytical & analysisType == 1:
        xL, L1, L2, Lsimple = getAnalyticalCurve()
        print("++++++++++++++++++++++++++++")
        print(xL)
        sns.lineplot(x=xL, y=Lsimple, color="red",
                     label="Analytical (simple)", linestyle="dashed")
        sns.lineplot(x=xL, y=(L1+L2)/2., color="red", label="Analytical")
        sns.lineplot(x=xL, y=L1, color="red", linestyle="dotted")
        sns.lineplot(x=xL, y=L2, color="red", linestyle="dotted")
    if analysisType == 2:
        plt.yscale('log')
        plt.ylim(ylims)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.xlim(xlims)
    plt.legend()
    plt.savefig(fileSaveFig, format='png')
    plt.clf()


def evaluate2(analysisType):
    if analysisType == 1:
        filenamedata = folderdata+"tips_"
        ylabel = "Number of tips"
        fileSaveFig = folder+'tips.png'
        column = 1
    elif analysisType == 2:
        filenamedata = folderdata+"orphantips_"
        ylabel = "Orphanage rate"
        fileSaveFig = folder+'orphanage.png'
        column = 2

    X = loadColumn(folderdata+"params", 0, 0)

    fig, ax = plt.subplots()
    for i in np.arange(len(X)):
        x = X[i]
        y_data = loadColumn(filenamedata+str(i), column, 2)

        # Some layout stuff ----------------------------------------------
        # Background color
        fig.patch.set_facecolor(BG_WHITE)
        ax.set_facecolor(BG_WHITE)
        # violins = ax.violinplot(
        #     y_data,
        #     positions=[[x]],
        #     widths=0.45,
        #     bw_method="silverman",
        #     showmeans=False,
        #     showmedians=False,
        #     showextrema=False
        # )
        # # Customize violins (remove fill, customize line, etc.)
        # for pc in violins["bodies"]:
        #     pc.set_facecolor("none")
        #     pc.set_edgecolor(GREY_LIGHT)
        #     pc.set_linewidth(1.4)
        #     pc.set_alpha(1)

        # Add boxplots ---------------------------------------------------
        # Note that properties about the median and the box are passed
        # as dictionaries.

        medianprops = dict(
            linewidth=4,
            color=GREY_DARK,
            solid_capstyle="butt"
        )
        boxprops = dict(
            linewidth=2,
            color=GREY_DARK
        )
        c = GREY_LIGHT
        bp = ax.boxplot(
            y_data,
            widths=max(xlims)/len(X)*.8,
            positions=[x],
            showfliers=False,  # Do not show the outliers beyond the caps.
            showcaps=False,   # Do not show the caps
            # medianprops=medianprops,
            # whiskerprops=boxprops,
            # boxprops=boxprops
            # fill the boxplot with color
            patch_artist=True,
            boxprops=dict(facecolor=c, color=c),
            capprops=dict(color=c),
            whiskerprops=dict(color=c),
            flierprops=dict(color=c, markeredgecolor=c),
            medianprops=dict(color=c),
        )
        for element in ['boxes', 'whiskers', 'fliers', 'means', 'medians', 'caps']:
            plt.setp(bp[element], color=BLACK)

        for patch in bp['boxes']:
            patch.set(facecolor=c)
        # use seaborn instead - the problem is that the position is hard to define with x axis
        # data = np.concatenate([[y_data, np.ones(len(y_data))*x]], axis=1)
        # df = pd.DataFrame(columns=['value', 'site'], data=data.T)
        # df['value'] = df['value'].astype(float)
        # sns.boxplot(x='site', y='value',  data=df)

    if analysisType == 2:
        plt.yscale('log')
        plt.ylim(ylims)
    plt.xlim(xlims)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.savefig(fileSaveFig+'_v2.png', format='png')
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


def getAnalyticalCurve():
    lam = 100.
    h = 1.
    k = 2.
    q = .25
    D = 6.

    # load X data
    Xdata = loadColumn(folderdata+"params", 0, 0)
    X = np.arange(100)/100.
    # X = max(Xdata)*np.ones(len(X))-(max(Xdata)-min(Xdata))*X
    X = 20.*X+1.05
    D = X-1

    p = 1-q

    # solve the following numerically
    # L0 = p*k(h-D*np.exp(-D*p*k/L0))/(p*k-1+np.exp(-D*p*k/L0))*lam*h
    L0 = X*0+1
    print(L0[:3])
    for i in range(1001):
        L1 = L0
        L0 = p*k*(h-D*np.exp(-D*p*k/L0))/(p*k-1+np.exp(-D*p*k/L0))
        delta = L0-L1
        pos = 24
        print("x,L0,delta = ", X[pos], L0[pos], np.log(abs(delta[pos])))
    L2 = L0*lam
    L1 = L1*lam
    L1[L1 < 0] = np.NaN
    L2[L1 < 0] = np.NaN

    print("------- Numerical precision ---- ")
    print(len(L1), len(delta), len(X))
    input = X[L0 > 0]
    result = L1[L0 > 0]
    error = abs(delta[L0 > 0])
    print("Xvalue\tResult\tError")
    for i in range(len(input)):
        print(input[i], "\t", result[i], "\t", error[i])

    # simple equation instead
    Lsimple = p*k/(p*k-1)*h*lam*np.ones(len(X))
    Lsimple[Lsimple < 0] = np.NaN

    return X, L1, L2, Lsimple


# needs to be at the very end of the file
if __name__ == '__main__':
    main()
